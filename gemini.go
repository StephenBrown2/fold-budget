package main

import (
	"time"
)

type GeminiCard struct {
	Reference   int        `csv:"Reference Number"`
	PostDate    geminiDate `csv:"Transaction Post Date"`
	Description string     `csv:"Description of Transaction"`
	Type        string     `csv:"Transaction Type"`
	Amount      float64    `csv:"Amount"`
}

type geminiDate struct {
	time.Time
}

func (d *geminiDate) UnmarshalCSV(data []byte) (err error) {
	d.Time, err = time.Parse("01/02/06", string(data))
	return err
}

func (d *geminiDate) MarshalCSV() ([]byte, error) {
	return []byte(d.Format("01/02/06")), nil
}

func (d *geminiDate) String() string {
	return d.Format(time.DateOnly)
}
