package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jszwec/csvutil"
)

var (
	dryRun    bool
	isBitcoin bool
)

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run, don't write to file")
	flag.BoolVar(&isBitcoin, "bitcoin", false, "Bitcoin account")
	flag.Parse()
}

func getFilename(data []YNAB, isBitcoin bool) string {
	oldestDate := time.Now().Add(time.Hour * 24 * 365)
	newestDate := time.Time{}
	// Iterate over the records to find the oldest and newest dates
	for _, y := range data {
		if y.Date.Time.Before(oldestDate) {
			oldestDate = y.Date.Time
		}
		if y.Date.Time.After(newestDate) {
			newestDate = y.Date.Time
		}
	}
	if newestDate.IsZero() {
		fmt.Println("Newest date is zero, using current date")
		newestDate = time.Now()
	}
	if oldestDate.After(time.Now()) {
		fmt.Println("Oldest date is not found, using current date")
		oldestDate = time.Now()
	}
	outFileType := "checking"
	if isBitcoin {
		outFileType = "bitcoin"
	}
	if oldestDate.Equal(newestDate) {
		fmt.Println("Oldest and newest dates are the same, using only one date for filename")
		return fmt.Sprintf("fold_%s_%s.csv", outFileType, oldestDate.Format(time.DateOnly))
	}
	return fmt.Sprintf("fold_%s_%s_%s.csv", outFileType, oldestDate.Format(time.DateOnly), newestDate.Format(time.DateOnly))
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

	out := []YNAB{}

	if isBitcoin {
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

			y, e := record.ToYNAB()
			if e != nil {
				fmt.Printf("Error converting to YNAB: %v\n", e)
				continue
			}
			out = append(out, y)
		}
	} else {
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

			date := record.SettlementDate.Time.Local()
			payee := record.Description
			out = append(out, YNAB{
				Date:   ynabDate{date},
				Payee:  payee,
				Amount: record.Amount,
			})
		}
	}

	if !dryRun {
		// Create a new CSV file to write the output
		outFileName := getFilename(out, isBitcoin)
		outFile, err := os.Create(outFileName)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer outFile.Close()

		// Write the records to the new CSV file using gocsv
		b, err := csvutil.Marshal(out)
		if err != nil {
			fmt.Println("error:", err)
		}
		if _, err := outFile.Write(b); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Printf("\nOutput written to %s\n", outFileName)
	}
}
