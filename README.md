# go-money

`go-money` provides a deterministic, currency-aware Money type for financial applications.
All calculations are performed by an internal decimal engine based on `github.com/govalues/decimal`.
Decimal is not part of the public API and is never exposed to callers.

## Example

```go
price := money.New(19990, money.Currency{Code: "TRY", Scale: 2, Symbol: "₺"})

invoice, err := price.
	SubtractPercent(10).
	AddPercent(18)
if err != nil {
	log.Fatal(err)
}

fmt.Println(invoice.String())
// ₺212.29
```

## Formatting

Global format can be configured with `SetFormat`, and any amount can be formatted locally (per-call) with `Format`.

```go
_ = money.SetFormat(money.FormatConfig{
	DecimalSeparator:   ",",
	ThousandsSeparator: ".",
	SymbolPosition:     money.SymbolSuffix,
	SymbolKind:         money.SymbolUseCurrencySymbol,
	Space:              true,
})

fmt.Println(invoice.String())
// 212,29 ₺

text, _ := invoice.Format(money.FormatConfig{
	DecimalSeparator:   ",",
	ThousandsSeparator: " ",
	SymbolPosition:     money.SymbolSuffix,
	SymbolKind:         money.SymbolUseCurrencyCode,
	Space:              true,
})
fmt.Println(text)
// 212,29 TRY

eur := money.New(123456789, money.Currency{Code: "EUR", Scale: 2, Symbol: "€"})
text, _ = eur.Format(money.FormatConfig{
	DecimalSeparator:   ",",
	ThousandsSeparator: ".",
	SymbolPosition:     money.SymbolSuffix,
	SymbolKind:         money.SymbolUseCurrencySymbol,
	Space:              true,
})
fmt.Println(text)
// 1.234.567,89 €

jpy := money.New(123, money.Currency{Code: "JPY", Scale: 0, Symbol: "¥"})
fmt.Println(jpy.String())
// ¥123
```

## Notes

- Money stores values as int64 minor units with an attached currency.
- Operations are deterministic and error-driven.
- No floats and no decimal types in the public API; formatting is explicit via config.
