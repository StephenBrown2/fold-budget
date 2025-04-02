package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/jszwec/csvutil"
)

func skipToHeader(file *os.File, headerStruct any) (*csv.Reader, []string) {
	csvHeader, err := csvutil.Header(headerStruct, "csv")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = len(csvHeader)

	fmt.Printf("Checking for header: %q\n", csvHeader)
	for {
		record, readErr := reader.Read()
		if readErr != nil {
			if !errors.Is(readErr, csv.ErrFieldCount) {
				return reader, nil
			}
			continue
		}
		if slices.Equal(record, csvHeader) {
			break
		}
	}

	if _, seekErr := file.Seek(reader.InputOffset(), 0); seekErr != nil {
		fmt.Println("Error seeking to start of CSV:", seekErr)
		return reader, nil
	}
	return csv.NewReader(file), csvHeader
}

func writeCSV(outFileName string, outData any) (err error) {
	outFile := os.Stdout
	if dryRun {
		fmt.Println("Dry run, not writing to file.")
		fmt.Printf("Would write to %s\n", outFileName)
	} else {
		outFile, err = os.Create(outFileName)
		if err != nil {
			return fmt.Errorf("error creating file: %w", err)
		}
		defer outFile.Close()
	}
	b, err := csvutil.Marshal(outData)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	if _, err := outFile.Write(b); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}
	if !dryRun {
		fmt.Printf("\nOutput written to %s\n", outFileName)
	}
	return nil
}
