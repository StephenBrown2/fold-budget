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

	fmt.Printf("Checking for header:\n%q\n", csvHeader)
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
