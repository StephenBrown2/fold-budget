package main

import "time"

type YNAB struct {
	Date   ynabDate
	Payee  string
	Memo   string
	Amount string
}

type ynabDate struct {
	time.Time
}

func (d *ynabDate) UnmarshalCSV(data []byte) error {
	t, err := time.Parse(time.DateOnly, string(data))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d *ynabDate) MarshalCSV() ([]byte, error) {
	return []byte(d.Format(time.DateOnly)), nil
}

func (d *ynabDate) String() string {
	return d.Format(time.DateOnly)
}
