package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
)

var (
	dryRun        bool
	since         dateValue
	unit          = newEnumFlag([]string{"usd", "btc", "sats"}, "usd")
	inFormat      = newEnumFlag([]string{"bitcoin", "checking", "debit"}, "checking")
	budgetFormats = []string{"ynab", "lunchmoney"}
	taxFormats    = []string{"coinledger", "cointracker", "koinly", "irr"}
	outFormat     = newEnumFlag(slices.Concat(budgetFormats, taxFormats), "ynab")
	oldestDate    time.Time
	newestDate    time.Time
)

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run, don't write to file")
	flag.Var(inFormat, "from", inFormat.Usage("Input format"))
	flag.Var(outFormat, "to", outFormat.Usage("Output format"))
	flag.Var(&since, "since", "Include transactions since this date")
	flag.Var(unit, "unit", unit.Usage("Output currency unit"))

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nUsage: %s [flags] <csv-file>\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nNOTE:\n")
		fmt.Fprintf(w, "  The \"debit\" input format is an alias for \"checking\".\n\n")
		fmt.Fprintf(w, "  The following output formats are available for all input formats:\n")
		fmt.Fprintf(w, "\t%s\n\n", strings.Join(budgetFormats, ", "))
		fmt.Fprintf(w, "  The following output formats are only available for bitcoin CSVs:\n")
		fmt.Fprintf(w, "\t%s\n\n", strings.Join(taxFormats, ", "))
	}

	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("No CSV file specified.")
		flag.Usage()
		os.Exit(1)
	}

	if inFormat.String() == "" {
		fmt.Println("Input format not specified, defaulting to checking.")
		_ = inFormat.Set("checking")
	}
	if outFormat.String() == "" {
		fmt.Println("Output format is required. Use -to flag to specify.")
		flag.Usage()
		os.Exit(1)
	}

	if slices.Contains(taxFormats, outFormat.String()) && inFormat.String() != "bitcoin" {
		fmt.Println("Tax Ledger format is not supported for non-bitcoin accounts")
		os.Exit(1)
	}

	if unit.String() != "usd" && inFormat.String() != "bitcoin" {
		fmt.Println("Bitcoin output unit is not supported for non-bitcoin accounts")
		os.Exit(1)
	}

	if outFormat.String() == "ynab" && unit.String() == "btc" {
		fmt.Println("YNAB format does not support BTC output unit")
		os.Exit(1)
	}

	if outFormat.String() == "lunchmoney" && unit.String() == "sats" {
		fmt.Println("Lunch Money format does not support sats output unit")
		os.Exit(1)
	}
}

func main() {
	// Get the CSV file from the first argument
	csvFile := flag.Arg(0)
	// Open the CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	oldestDate = time.Now().Add(time.Hour * 24 * 365)
	newestDate = time.Time{}
	btctxns := []FoldBitcoin{}
	txns := []Transaction{}

	switch inFormat.String() {
	case "bitcoin":
		csvReader, csvHeader := skipToHeader(file, FoldBitcoin{})
		dec, err := csvutil.NewDecoder(csvReader, csvHeader...)
		if errors.Is(err, io.EOF) {
			fmt.Println("Reached end of file early. Is it a bitcoin statement?")
			return
		} else if err != nil {
			log.Fatal(err)
		}
		for {
			var record FoldBitcoin
			if err := dec.Decode(&record); errors.Is(err, io.EOF) {
				break
			} else if errors.Is(err, csv.ErrFieldCount) {
				fmt.Println("Skipping", err.Error())
				continue
			} else if err != nil {
				fmt.Printf("Error decoding: %v\n", err)
				continue
			}

			if record.DateUTC.Before(since.Time) {
				continue
			}
			if record.DateUTC.Before(oldestDate) {
				oldestDate = record.DateUTC.Time
			}
			if record.DateUTC.After(newestDate) {
				newestDate = record.DateUTC.Time
			}

			switch outFormat.String() {
			case "ynab", "lunchmoney":
				t, e := record.Transaction(unit.String())
				if e != nil {
					fmt.Printf("Error converting to budget transaction: %v\n", e)
					continue
				}
				txns = append(txns, t)
			case "coinledger", "cointracker", "koinly", "irr":
				btctxns = append(btctxns, record)
			}
		}
	case "checking", "debit":
		csvReader, csvHeader := skipToHeader(file, FoldCard{})
		dec, err := csvutil.NewDecoder(csvReader, csvHeader...)
		if errors.Is(err, io.EOF) {
			fmt.Println("Reached end of file early. Is it a debit card statement?")
			return
		} else if err != nil {
			log.Fatal(err)
		}
		for {
			var record FoldCard
			if err := dec.Decode(&record); errors.Is(err, io.EOF) {
				break
			} else if errors.Is(err, csv.ErrFieldCount) {
				fmt.Println("Skipping", err.Error())
				continue
			} else if err != nil {
				fmt.Printf("Error decoding: %v\n", err)
				continue
			}

			if record.SettlementDate.Before(since.Time) {
				continue
			}
			if record.SettlementDate.Before(oldestDate) {
				oldestDate = record.SettlementDate.Time
			}
			if record.SettlementDate.After(newestDate) {
				newestDate = record.SettlementDate.Time
			}

			date := record.SettlementDate.Time.Local()
			payee := record.Description
			txns = append(txns, Transaction{
				Date:   date,
				Payee:  payee,
				Amount: record.Amount,
			})
		}
	}

	// Create a new CSV file to write the output
	outFileName := getFilename(oldestDate, newestDate, inFormat, outFormat)
	switch outFormat.String() {
	case "irr":
		fmt.Println("Processing IRR...")
		costbasis, totalbtc := 0.0, 0.0
		for _, t := range btctxns {
			if t.PricePerCoinUSD.float64 != 0 {
				costbasis += t.AmountBTC * t.PricePerCoinUSD.float64
				totalbtc += t.AmountBTC
			}
		}
		fmt.Printf("Total cost basis: $%.2f\n", costbasis)
		fmt.Printf("Total BTC: %.8f\n", totalbtc)
		currentPrice, err := getCurrentPrice()
		if err != nil {
			fmt.Printf("Error getting current price: %v\n", err)
			return
		}
		currentvalue := totalbtc * currentPrice
		profit := currentvalue - costbasis
		fmt.Printf("Current price: $%.2f\n", currentPrice)
		fmt.Printf("Total USD: $%.2f\n", currentvalue)
		fmt.Printf("Total profit: $%.2f\n", profit)
		fmt.Printf("Total profit percentage: %.2f%%\n", profit/costbasis*100)
	case "coinledger":
		fmt.Println("Processing with CoinLedger format...")
		outData := []CoinLedger{}
		for _, t := range btctxns {
			outData = append(outData, t.ToCoinLedger())
		}

		if err := writeCSV(outFileName, outData); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	case "cointracker":
		fmt.Println("Processing with CoinTracker format...")
		outData := []CoinTracker{}
		for _, t := range btctxns {
			outData = append(outData, t.ToCoinTracker())
		}

		if err := writeCSV(outFileName, outData); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	case "koinly":
		fmt.Println("Processing with Koinly format...")
		outData := []Koinly{}
		for _, t := range btctxns {
			outData = append(outData, t.ToKoinly())
		}

		if err := writeCSV(outFileName, outData); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	case "lunchmoney":
		fmt.Println("Processing with Lunch Money format...")
		outData := []LunchMoney{}
		for _, t := range txns {
			var amt string
			switch unit.String() {
			case "btc":
				amt = fmt.Sprintf("%.8f", t.Amount)
			case "sats":
				amt = fmt.Sprintf("%.0f", t.Amount)
			case "usd":
				amt = fmt.Sprintf("%.2f", t.Amount)
			}
			outData = append(outData, LunchMoney{
				Date:       lmDate{t.Date},
				Payee:      t.Payee,
				Notes:      t.Memo,
				Amount:     amt,
				Categories: t.Categories,
				Tags:       t.Tags,
			})
		}

		if err := writeCSV(outFileName, outData); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	case "ynab":
		fmt.Println("Processing with YNAB format...")
		outData := []YNAB{}
		for _, t := range txns {
			var amt string
			switch unit.String() {
			case "btc":
				amt = fmt.Sprintf("%.8f", t.Amount)
			case "sats":
				amt = fmt.Sprintf("%.0f", t.Amount)
			case "usd":
				amt = fmt.Sprintf("%.2f", t.Amount)
			}
			outData = append(outData, YNAB{
				Date:   ynabDate{t.Date},
				Payee:  t.Payee,
				Memo:   t.Memo,
				Amount: amt,
			})
		}

		if err := writeCSV(outFileName, outData); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
}
