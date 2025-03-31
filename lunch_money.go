package main

import "time"

type LunchMoney struct {
	Date       lmDate   `csv:"date"`
	Payee      string   `csv:"payee"`
	Notes      string   `csv:"notes"`
	Amount     float64  `csv:"amount"`
	Categories []string `csv:"categories,omitempty"`
	Tags       []string `csv:"tags,omitempty"`
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
	return []byte(d.Time.Format(time.DateOnly)), nil
}

func (d *lmDate) String() string {
	return d.Time.Format(time.DateOnly)
}
