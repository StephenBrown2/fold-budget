package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type FoldBitcoin struct {
	ReferenceID     string   `csv:"Reference ID"`
	DateUTC         foldUTC  `csv:"Date (UTC)"`
	TransactionType string   `csv:"Transaction Type"`
	Description     string   `csv:"Description"`
	Asset           string   `csv:"Asset"`
	AmountBTC       float64  `csv:"Amount (BTC)"`
	PricePerCoinUSD optFloat `csv:"Price per Coin (USD)"`
	SubtotalUSD     optFloat `csv:"Subtotal (USD)"`
	FeeUSD          optFloat `csv:"Fee (USD)"`
	TotalUSD        optFloat `csv:"Total (USD)"`
	TransactionID   string   `csv:"Transaction ID"`
}

type optFloat struct {
	float64
}

func (o *optFloat) UnmarshalCSV(data []byte) (err error) {
	if string(data) == "" {
		o.float64 = 0
		return nil
	}
	o.float64, err = strconv.ParseFloat(string(data), 64)
	return err
}

type foldUTC struct {
	time.Time
}

func (d *foldUTC) UnmarshalCSV(data []byte) (err error) {
	d.Time, err = time.Parse("2006-01-02 15:04:05.999999-07:00", string(data))
	return err
}

func (d *foldUTC) MarshalCSV() ([]byte, error) {
	return []byte(d.Time.Format(time.RFC3339)), nil
}

func (d *foldUTC) String() string {
	return d.Time.Format(time.RFC3339)
}

func (record FoldBitcoin) ToYNAB() (YNAB, error) {
	date := record.DateUTC.Time.Local()

	payee := record.Description
	switch payee {
	case "Direct to Bitcoin Purchase":
		payee = "Fold Direct to Bitcoin Purchase"
	case "Push to Card":
		payee = "Fold Push to Card"
	case "Purchase":
		payee = "Fold Bitcoin Purchase"
	case "Auto-Stack Purchase":
		payee = "Fold Auto-Stack Bitcoin Purchase"
	case "Receive":
		payee = "Fold Receive Bitcoin"
	}

	amount := record.SubtotalUSD.float64 * -1
	price := record.PricePerCoinUSD.float64
	if record.TotalUSD.float64 == 0 {
		var priceErr error
		price, priceErr = getHistoricalPrice(record.DateUTC.Time)
		if priceErr != nil {
			return YNAB{}, fmt.Errorf("Error getting historical price: %w", priceErr)
		}
		amount = math.Round(record.AmountBTC*price*100) / 100
	}

	memo := strings.Builder{}
	memo.WriteString(fmt.Sprintf("Cost Basis: %.2f ⋅ %.8f BTC (FX rate: %.2f)", amount, record.AmountBTC, price))
	if record.TransactionID != "" {
		memo.WriteString(fmt.Sprintf(" Transaction ID: %s", record.TransactionID))
	}
	memo.WriteString(fmt.Sprintf(" @ %s", date.Format(time.RFC822)))

	return YNAB{
		Date:   ynabDate{date},
		Payee:  payee,
		Memo:   memo.String(),
		Amount: amount,
	}, nil
}
