package main

import "time"

type CoinTracker struct {
	Date             coinTrackerDate `csv:"Date"`              // Date of transaction
	ReceivedQuantity string          `csv:"Received Quantity"` // Amount of crypto or cash received
	ReceivedCurrency string          `csv:"Received Currency"` // Type of crypto received
	SentQuantity     string          `csv:"Sent Quantity"`     // Amount of crypto or cash sent
	SentCurrency     string          `csv:"Sent Currency"`     // Type of crypto sent
	FeeAmount        float64         `csv:"Fee Amount"`        // Transaction fee amount in the currency it was paid
	FeeCurrency      string          `csv:"Fee Currency"`      // Type of currency that your transaction fee was paid in
	Tag              CoinTrackerTag  `csv:"Tag"`               // [CoinTracker CSV tags](https://support.cointracker.io/hc/en-us/articles/4413049710225): Use tags to categorize send/receive transactions by type for better tracking and tax purposes.* Do not use tags for trades or transfers.
}

type coinTrackerDate struct {
	time.Time
}

func (d *coinTrackerDate) UnmarshalCSV(data []byte) (err error) {
	d.Time, err = time.Parse("01/02/2006 15:04:05", string(data))
	return err
}

func (d *coinTrackerDate) MarshalCSV() ([]byte, error) {
	return []byte(d.Time.UTC().Format("01/02/2006 15:04:05")), nil
}

func (d *coinTrackerDate) String() string {
	return d.Time.UTC().Format("01/02/2006 15:04:05")
}

type CoinTrackerTag string

// Transaction categories label cryptocurrency transactions for tax and reporting purposes. You can manually adjust the category of automatically synced transactions as needed.
