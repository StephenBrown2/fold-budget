package main

import "time"

type CoinLedgerType string

const (
	CoinLedgerAirdrop         CoinLedgerType = "Airdrop"
	CoinLedgerCasualtyLoss    CoinLedgerType = "Casualty Loss"
	CoinLedgerDeposit         CoinLedgerType = "Deposit"
	CoinLedgerGiftReceived    CoinLedgerType = "Gift Received"
	CoinLedgerGiftSent        CoinLedgerType = "Gift Sent"
	CoinLedgerHardFork        CoinLedgerType = "Hard Fork"
	CoinLedgerIncome          CoinLedgerType = "Income"
	CoinLedgerInterest        CoinLedgerType = "Interest"
	CoinLedgerInterestPayment CoinLedgerType = "Interest Payment"
	CoinLedgerInvestmentLoss  CoinLedgerType = "Investment Loss"
	CoinLedgerMerchantPayment CoinLedgerType = "Merchant Payment"
	CoinLedgerMining          CoinLedgerType = "Mining"
	CoinLedgerStaking         CoinLedgerType = "Staking"
	CoinLedgerTheftLoss       CoinLedgerType = "Theft Loss"
	CoinLedgerTrade           CoinLedgerType = "Trade"
	CoinLedgerWithdrawal      CoinLedgerType = "Withdrawal"
)

type CoinLedger struct {
	DateUTC        clUTC          `csv:"Date (UTC)"`
	Platform       string         `csv:"Platform (Optional),omitempty"`
	AssetSent      string         `csv:"Asset Sent"`
	AmountSent     float64        `csv:"Amount Sent"`
	AssetReceived  string         `csv:"Asset Received"`
	AmountReceived string         `csv:"Amount Received"` // Note: this is a string to ensure proper formatting
	FeeCurrency    string         `csv:"Fee Currency (Optional),omitempty"`
	FeeAmount      float64        `csv:"Fee Amount (Optional),omitempty"`
	Type           CoinLedgerType `csv:"Type"`
	Description    string         `csv:"Description (Optional),omitempty"`
	TxHash         string         `csv:"TxHash (Optional),omitempty"`
}

type clUTC struct {
	time.Time
}

func (d *clUTC) UnmarshalCSV(data []byte) (err error) {
	d.Time, err = time.Parse("01-02-2006 15:04:05", string(data))
	return err
}

func (d *clUTC) MarshalCSV() ([]byte, error) {
	return []byte(d.Time.Format("01-02-2006 15:04:05")), nil
}

func (d *clUTC) String() string {
	return d.Time.Format("01-02-2006 15:04:05")
}
