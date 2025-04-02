# fold-budget

`fold-budget` is a CLI tool to work with your [Fold](https://foldapp.com/)
card or bitcoin transaction CSV exports, and translate them to your budget,
currently:

- YNAB
- LunchMoney

And for bitcoin account transactions, the following additional targets are
supported:

- CoinLedger
- CoinTracker
- Koinly

To install, make sure you have [Go](https://go.dev/dl/) installed, then run:

```sh
go install https://github.com/StephenBrown2/fold-budget@latest
```

Then, you should be able to run it:

```sh
$ fold-budget -help
Usage: fold-budget [flags] <csv-file>
  -dry-run
        Dry run, don't write to file
  -from value
        Input format (bitcoin or checking, default: checking)
  -since value
        Include transactions since this date
  -to value
        Output format (one of: ynab, lunchmoney, coinledger, cointracker, koinly)
```

Once you have a downloaded CSV file, you can choose your input and output
formats, and `fold-budget` will write a new file based on those and the start
and end dates for the included transactions:

```sh
$ fold-budget -from=bitcoin -to=koinly ~/Downloads/fold-bitcoin-transaction-history-2025-03-31.csv
Checking for header: ["Reference ID" "Date (UTC)" "Transaction Type" "Description" "Asset" "Amount (BTC)" "Price per Coin (USD)" "Subtotal (USD)" "Fee (USD)" "Total (USD)" "Transaction ID"]
Getting historical price for: Sun, 02 Feb 2025 21:40:43 EST
Skipping record on line 27: wrong number of fields
Skipping record on line 28: wrong number of fields
Processing with Koinly format...
Getting historical price for: Sun, 02 Feb 2025 21:40:43 EST

Output written to fold_bitcoin_to_koinly_2025-01-09_2025-03-31.csv
```

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. For more information, please refer to <https://unlicense.org>.
