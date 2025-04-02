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

func (record *FoldBitcoin) USDPerCoin() (float64, error) {
	price := record.PricePerCoinUSD.float64
	if price == 0 {
		var priceErr error
		price, priceErr = getHistoricalPrice(record.DateUTC.Time)
		if priceErr != nil {
			return 0, fmt.Errorf("error getting historical price: %w", priceErr)
		}
		record.PricePerCoinUSD = optFloat{price}
	}
	return price, nil
}

func (record FoldBitcoin) Transaction() (Transaction, error) {
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
	price, priceErr := record.USDPerCoin()
	if priceErr != nil {
		return Transaction{}, fmt.Errorf("error getting historical price: %w", priceErr)
	}
	if record.TotalUSD.float64 == 0 {
		amount = math.Round(record.AmountBTC*price*100) / 100
	}

	memo := strings.Builder{}
	memo.WriteString(fmt.Sprintf("Cost Basis: %.2f ⋅ %.8f BTC (FX rate: %.2f)", amount, record.AmountBTC, price))
	if record.TransactionID != "" {
		memo.WriteString(fmt.Sprintf(" Transaction ID: %s", record.TransactionID))
	}
	memo.WriteString(fmt.Sprintf(" @ %s", date.Format(time.RFC822)))

	return Transaction{
		Date:   date,
		Payee:  payee,
		Memo:   memo.String(),
		Amount: amount,
	}, nil
}

func (record FoldBitcoin) ToCoinLedger() CoinLedger {
	sent := math.Abs(record.SubtotalUSD.float64)
	received := math.Abs(record.AmountBTC)
	assetSent := "USD"
	assetReceived := "BTC"

	var txType CoinLedgerTag
	switch record.TransactionType {
	case "Purchase", "Sale":
		txType = CoinLedgerTrade
	case "Deposit":
		txType = CoinLedgerDeposit
	case "Withdrawal":
		txType = CoinLedgerWithdrawal
		sent = math.Abs(record.AmountBTC)
		received = math.Abs(record.SubtotalUSD.float64)
		assetSent = "BTC"
		assetReceived = "USD"
	default:
		fmt.Println("Unknown transaction type:", record.TransactionType)
		if record.AmountBTC > 0 {
			fmt.Println("Assuming deposit for positive BTC amount")
			txType = CoinLedgerDeposit
		} else {
			fmt.Println("Assuming withdrawal for negative BTC amount")
			txType = CoinLedgerWithdrawal
		}
	}

	price, priceErr := record.USDPerCoin()
	if priceErr != nil {
		fmt.Printf("error getting historical BTC price: %s\n", priceErr)
	}

	var amountReceived, amountSent string
	switch assetReceived {
	case "BTC":
		amountReceived = fmt.Sprintf("%0.8f", received)
		amountSent = fmt.Sprintf("%0.2f", sent)
	case "USD":
		amountReceived = fmt.Sprintf("%0.2f", received)
		amountSent = fmt.Sprintf("%0.8f", sent)
	}

	return CoinLedger{
		DateUTC:        coinLedgerDate{record.DateUTC.Time},
		Platform:       "Fold",
		AssetSent:      assetSent,
		AmountSent:     amountSent,
		AssetReceived:  assetReceived,
		AmountReceived: amountReceived,
		FeeCurrency:    "USD",
		FeeAmount:      record.FeeUSD.float64,
		Type:           txType,
		Description:    fmt.Sprintf("%s, FX rate: %.2f", record.Description, price),
		TxHash:         record.TransactionID,
	}
}

func (record FoldBitcoin) ToCoinTracker() CoinTracker {
	sent := math.Abs(record.SubtotalUSD.float64)
	received := math.Abs(record.AmountBTC)
	assetSent := "USD"
	assetReceived := "BTC"

	switch record.TransactionType {
	case "Withdrawal", "Sale":
		sent = math.Abs(record.AmountBTC)
		received = math.Abs(record.SubtotalUSD.float64)
		assetSent = "BTC"
		assetReceived = "USD"
	}

	var amountReceived, amountSent string
	switch assetReceived {
	case "BTC":
		amountReceived = fmt.Sprintf("%0.8f", received)
		amountSent = fmt.Sprintf("%0.2f", sent)
	case "USD":
		amountReceived = fmt.Sprintf("%0.2f", received)
		amountSent = fmt.Sprintf("%0.8f", sent)
	}

	return CoinTracker{
		Date:             coinTrackerDate{record.DateUTC.Time},
		ReceivedQuantity: amountReceived,
		ReceivedCurrency: assetReceived,
		SentQuantity:     amountSent,
		SentCurrency:     assetSent,
		FeeCurrency:      "USD",
		FeeAmount:        record.FeeUSD.float64,
	}
}

func (record FoldBitcoin) ToKoinly() Koinly {
	sent := math.Abs(record.SubtotalUSD.float64)
	received := math.Abs(record.AmountBTC)
	assetSent := "USD"
	assetReceived := "BTC"

	switch record.TransactionType {
	case "Withdrawal", "Sale":
		sent = math.Abs(record.AmountBTC)
		received = math.Abs(record.SubtotalUSD.float64)
		assetSent = "BTC"
		assetReceived = "USD"
	}

	price, priceErr := record.USDPerCoin()
	if priceErr != nil {
		fmt.Printf("error getting historical BTC price: %s\n", priceErr)
	}

	var networth float64
	var amountReceived, amountSent string
	switch assetReceived {
	case "BTC":
		amountReceived = fmt.Sprintf("%0.8f", received)
		amountSent = fmt.Sprintf("%0.2f", sent)
		networth = price * math.Abs(record.AmountBTC)
	case "USD":
		amountReceived = fmt.Sprintf("%0.2f", received)
		amountSent = fmt.Sprintf("%0.8f", sent)
		networth = math.Abs(record.SubtotalUSD.float64)
	}

	return Koinly{
		Date:             koinlyDate{record.DateUTC.Time},
		SentAmount:       amountSent,
		SentCurrency:     assetSent,
		ReceivedAmount:   amountReceived,
		ReceivedCurrency: assetReceived,
		FeeCurrency:      "USD",
		FeeAmount:        record.FeeUSD.float64,
		NetWorthAmount:   math.Round(networth*100) / 100,
		NetWorthCurrency: "USD",
		Description:      fmt.Sprintf("%s, FX rate: %.2f", record.Description, price),
		TxHash:           record.TransactionID,
	}
}
