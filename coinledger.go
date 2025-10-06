package main

import "time"

type CoinLedgerTag string

const (
	CoinLedgerAirdrop         CoinLedgerTag = "Airdrop"
	CoinLedgerCasualtyLoss    CoinLedgerTag = "Casualty Loss"
	CoinLedgerDeposit         CoinLedgerTag = "Deposit"
	CoinLedgerGiftReceived    CoinLedgerTag = "Gift Received"
	CoinLedgerGiftSent        CoinLedgerTag = "Gift Sent"
	CoinLedgerHardFork        CoinLedgerTag = "Hard Fork"
	CoinLedgerIncome          CoinLedgerTag = "Income"
	CoinLedgerInterest        CoinLedgerTag = "Interest"
	CoinLedgerInterestPayment CoinLedgerTag = "Interest Payment"
	CoinLedgerInvestmentLoss  CoinLedgerTag = "Investment Loss"
	CoinLedgerMerchantPayment CoinLedgerTag = "Merchant Payment"
	CoinLedgerMining          CoinLedgerTag = "Mining"
	CoinLedgerStaking         CoinLedgerTag = "Staking"
	CoinLedgerTheftLoss       CoinLedgerTag = "Theft Loss"
	CoinLedgerTrade           CoinLedgerTag = "Trade"
	CoinLedgerWithdrawal      CoinLedgerTag = "Withdrawal"
)

type CoinLedger struct {
	DateUTC        coinLedgerDate `csv:"Date (UTC)"`
	Platform       string         `csv:"Platform (Optional),omitempty"`
	AssetSent      string         `csv:"Asset Sent"`
	AmountSent     string         `csv:"Amount Sent"` // Note: this is a string to ensure proper formatting
	AssetReceived  string         `csv:"Asset Received"`
	AmountReceived string         `csv:"Amount Received"` // Note: this is a string to ensure proper formatting
	FeeCurrency    string         `csv:"Fee Currency (Optional),omitempty"`
	FeeAmount      float64        `csv:"Fee Amount (Optional),omitempty"`
	Type           CoinLedgerTag  `csv:"Type"`
	Description    string         `csv:"Description (Optional),omitempty"`
	TxHash         string         `csv:"TxHash (Optional),omitempty"`
}

type coinLedgerDate struct {
	time.Time
}

func (d *coinLedgerDate) UnmarshalCSV(data []byte) (err error) {
	d.Time, err = time.Parse("01-02-2006 15:04:05", string(data))
	return err
}

func (d *coinLedgerDate) MarshalCSV() ([]byte, error) {
	return []byte(d.Format("01-02-2006 15:04:05")), nil
}

func (d *coinLedgerDate) String() string {
	return d.Format("01-02-2006 15:04:05")
}
