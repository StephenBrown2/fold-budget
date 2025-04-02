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

	"github.com/jszwec/csvutil"
)

var (
	dryRun    bool
	since     dateValue
	inFormat  inputFormat
	outFormat outputFormat
)

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run, don't write to file")
	flag.Var(&inFormat, "from", "Input format (bitcoin or checking, default: checking)")
	flag.Var(&outFormat, "to", "Output format (one of: ynab, lunchmoney, coinledger, cointracker, koinly)")
	flag.Var(&since, "since", "Include transactions since this date")

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s [flags] <csv-file>\n", os.Args[0])
		flag.PrintDefaults()
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

	if slices.Contains([]string{"coinledger", "cointracker", "koinly"}, outFormat.String()) && inFormat.String() != "bitcoin" {
		fmt.Println("Tax Ledger format is not supported for non-bitcoin accounts")
		return
	}
}

func main() {
	// Get the CSV file from the first argument
	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a CSV file.")
		return
	}
	csvFile := flag.Arg(0)
	// Open the CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

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

			t, e := record.Transaction()
			if e != nil {
				fmt.Printf("Error converting to budget transaction: %v\n", e)
				continue
			}
			txns = append(txns, t)
			btctxns = append(btctxns, record)
		}
	case "checking", "card":
		csvReader, csvHeader := skipToHeader(file, FoldCard{})
		dec, err := csvutil.NewDecoder(csvReader, csvHeader...)
		if errors.Is(err, io.EOF) {
			fmt.Println("Reached end of file early. Is it a card statement?")
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
	outFileName := getFilename(txns, inFormat, outFormat)
	switch outFormat.String() {
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
			outData = append(outData, LunchMoney{
				Date:       lmDate{t.Date},
				Payee:      t.Payee,
				Notes:      t.Memo,
				Amount:     t.Amount,
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
			outData = append(outData, YNAB{
				Date:   ynabDate{t.Date},
				Payee:  t.Payee,
				Memo:   t.Memo,
				Amount: t.Amount,
			})
		}

		if err := writeCSV(outFileName, outData); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
}
