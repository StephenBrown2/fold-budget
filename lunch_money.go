package main

import (
	"strings"
	"time"
)

type LunchMoney struct {
	Date       lmDate  `csv:"date"`
	Payee      string  `csv:"payee"`
	Notes      string  `csv:"notes"`
	Amount     string  `csv:"amount"`
	Categories Strings `csv:"categories,omitempty"`
	Tags       Strings `csv:"tags,omitempty"`
}

type lmDate struct {
	time.Time
}

func (d *lmDate) UnmarshalCSV(data []byte) error {
	t, err := time.Parse(time.DateOnly, string(data))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d *lmDate) MarshalCSV() ([]byte, error) {
	return []byte(d.Format(time.DateOnly)), nil
}

func (d *lmDate) String() string {
	return d.Format(time.DateOnly)
}

type Strings []string

func (s Strings) UnmarshalCSV(data []byte) error {
	parts := strings.Split(string(data), ",")
	for _, part := range parts {
		s = append(s, strings.TrimSpace(part))
	}
	return nil
}

func (s Strings) MarshalCSV() ([]byte, error) {
	return []byte(strings.Join(s, ",")), nil // strings.Join takes []string but it will also accept Strings
}

func (s Strings) String() string {
	return strings.Join(s, ", ")
}
