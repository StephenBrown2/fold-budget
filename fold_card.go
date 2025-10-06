package main

import "time"

type FoldCard struct {
	TransactionID  string   `csv:"Transaction ID"`
	SettlementDate foldDate `csv:"Settlement Date"`
	Description    string   `csv:"Description"`
	Amount         float64  `csv:"Amount"`
}

type foldDate struct {
	time.Time
}

func (d *foldDate) UnmarshalCSV(data []byte) (err error) {
	d.Time, err = time.Parse("2006-01-02 15:04:05-07:00", string(data))
	return err
}

func (d *foldDate) MarshalCSV() ([]byte, error) {
	return []byte(d.Format(time.RFC3339)), nil
}

func (d *foldDate) String() string {
	return d.Format(time.RFC3339)
}
