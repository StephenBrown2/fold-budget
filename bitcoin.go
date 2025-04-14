package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func getHistoricalPrice(date time.Time) (float64, error) {
	fmt.Printf("Getting historical price for: %s\n", date.Local().Format(time.RFC1123))
	url := fmt.Sprintf("https://mempool.space/api/v1/historical-price?currency=USD&timestamp=%d", date.Unix())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	type Price struct {
		Timestamp int64   `json:"time"`
		USD       float64 `json:"USD"`
	}
	type Btc struct {
		Prices        []Price            `json:"prices"`
		ExchangeRates map[string]float64 `json:"exchangeRates"`
	}
	var btc Btc
	if err := json.NewDecoder(response.Body).Decode(&btc); err != nil {
		return 0, err
	}

	return btc.Prices[0].USD, nil
}

func getCurrentPrice() (float64, error) {
	fmt.Printf("Getting current price for: %s\n", time.Now().Local().Format(time.RFC1123))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://pricing.bitcoin.block.xyz/current-price", nil)
	if err != nil {
		return 0, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	type BlockPrice struct {
		Amount                         string `json:"amount"`
		LastUpdatedAtInUTCEpochSeconds string `json:"last_updated_at_in_utc_epoch_seconds"`
		Currency                       string `json:"currency"`
		Version                        string `json:"version"`
		Base                           string `json:"base"`
	}
	var btc BlockPrice
	if err := json.NewDecoder(response.Body).Decode(&btc); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(btc.Amount, 64)
}
