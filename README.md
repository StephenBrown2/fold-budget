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

## Installation

To install, make sure you have [Go](https://go.dev/dl/) installed, then run:

```sh
go install github.com/StephenBrown2/fold-budget@latest
```

## Usage

fold-budget has a traditional command-line interface:

```sh
fold-budget [flags] <csv-file>
```

### Flags

- `-dry-run`: Dry run, don't write to file
- `-from value`: Input format (bitcoin, checking, debit, gemini) - default: checking
- `-to value`: Output format (ynab, lunchmoney, coinledger, cointracker, koinly, irr) - default: ynab
- `-since value`: Include transactions since this date (YYYY-MM-DD format)
- `-unit value`: Output currency unit (usd, btc, sats) - default: usd
- `-tui`: Use interactive TUI mode

### Examples

Convert a bitcoin CSV to Koinly format:

```sh
fold-budget -from=bitcoin -to=koinly ~/Downloads/fold-bitcoin-transaction-history-2025-03-31.csv
```

Convert a checking account CSV to YNAB format:

```sh
fold-budget -from=checking -to=ynab ~/Downloads/statement.csv
```

## Supported Formats

### Input Formats

- **bitcoin**: Fold Bitcoin transaction history
- **checking/debit**: Fold checking account transactions
- **gemini**: Gemini card transactions

### Output Formats

**Available for all input formats:**

- **ynab**: You Need A Budget format
- **lunchmoney**: Lunch Money budgeting app format

**Available only for bitcoin CSVs:**

- **coinledger**: CoinLedger tax reporting format
- **cointracker**: CoinTracker tax reporting format
- **koinly**: Koinly tax reporting format
- **irr**: Internal Rate of Return calculation (console output only)

### Currency Units

- **usd**: US Dollars (available for all formats)
- **btc**: Bitcoin (bitcoin input only)
- **sats**: Satoshis (bitcoin input only, not supported by Lunch Money)

## Format Compatibility

| Output Format | USD | BTC | Sats |
|--------------|-----|-----|------|
| YNAB         | ✅  | ❌  | ✅   |
| Lunch Money  | ✅  | ✅  | ❌   |
| CoinLedger   | ✅  | ✅  | ✅   |
| CoinTracker  | ✅  | ✅  | ✅   |
| Koinly       | ✅  | ✅  | ✅   |

## Usage Examples

Bitcoin transactions to Koinly:

```sh
$ fold-budget -from=bitcoin -to=koinly ~/Downloads/fold-bitcoin-transaction-history-2025-03-31.csv
Checking for header: ["Reference ID" "Date (UTC)" "Transaction Type" "Description" "Asset" "Amount (BTC)" "Price per Coin (USD)" "Subtotal (USD)" "Fee (USD)" "Total (USD)" "Transaction ID"]
Processing with Koinly format...
Output written to fold_bitcoin_to_koinly_2025-01-10_2025-03-31.csv
```

Checking account transactions to YNAB:

```sh
$ fold-budget -from=checking -to=ynab ~/Downloads/statement.csv
Checking for header: ["Transaction ID" "Settlement Date" "Description" "Amount"]
Processing with YNAB format...
Output written to fold_checking_to_ynab_2025-03-01_2025-03-25.csv
```

## Notes

- The "debit" input format is an alias for "checking"
- Some rows may be skipped because Fold includes statement metadata along with transaction data
- When using bitcoin formats, the tool can fetch historical prices for transactions missing price data
- The TUI mode provides better error handling and validation than the CLI mode

## License

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. For more information, please refer to <https://unlicense.org>.
