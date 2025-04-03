package main

import "time"

type TaxRecord struct {
	DateTime       time.Time
	AssetSent      string
	AmountSent     string
	AssetReceived  string
	AmountReceived string
	FeeAsset       string
	FeeAmount      float64
	Description    string
	TxHash         string
}
