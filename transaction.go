package main

import (
	"fmt"
	"time"
)

type Transaction struct {
	Date       time.Time
	Payee      string
	Memo       string
	Amount     float64
	Categories []string
	Tags       []string
}

func getFilename(data []Transaction, from inputFormat, to outputFormat) string {
	oldestDate := time.Now().Add(time.Hour * 24 * 365)
	newestDate := time.Time{}
	// Iterate over the records to find the oldest and newest dates
	for _, t := range data {
		if t.Date.Before(oldestDate) {
			oldestDate = t.Date
		}
		if t.Date.After(newestDate) {
			newestDate = t.Date
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
	if oldestDate.Equal(newestDate) {
		fmt.Println("Oldest and newest dates are the same, using only one date for filename")
		return fmt.Sprintf("fold_%s_to_%s_%s.csv", from, to, oldestDate.Format(time.DateOnly))
	}
	return fmt.Sprintf("fold_%s_to_%s_%s_%s.csv", from, to, oldestDate.Format(time.DateOnly), newestDate.Format(time.DateOnly))
}
