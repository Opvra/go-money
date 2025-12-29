package money

import "testing"

func TestAddSub(t *testing.T) {
	usd := Currency{Code: "USD", Scale: 2, Symbol: "$"}
	a := New(1050, usd)
	b := New(250, usd)

	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("add error: %v", err)
	}
	if got := sum.Amount(); got != 1300 {
		t.Fatalf("sum amount = %d", got)
	}

	diff, err := a.Sub(b)
	if err != nil {
		t.Fatalf("sub error: %v", err)
	}
	if got := diff.Amount(); got != 800 {
		t.Fatalf("diff amount = %d", got)
	}
}

func TestCurrencyMismatch(t *testing.T) {
	usd := Currency{Code: "USD", Scale: 2, Symbol: "$"}
	eur := Currency{Code: "EUR", Scale: 2, Symbol: "€"}
	_, err := New(100, usd).Add(New(100, eur))
	if err != ErrCurrencyMismatch {
		t.Fatalf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestPercentExample(t *testing.T) {
	try := Currency{Code: "TRY", Scale: 2, Symbol: "₺"}
	price := New(19990, try)

	invoice, err := price.SubtractPercent(10)
	if err != nil {
		t.Fatalf("subtract percent error: %v", err)
	}
	invoice, err = invoice.AddPercent(18)
	if err != nil {
		t.Fatalf("add percent error: %v", err)
	}

	if got := invoice.String(); got != "₺212.29" {
		t.Fatalf("invoice string = %s", got)
	}
	if got := invoice.Amount(); got != 21229 {
		t.Fatalf("invoice amount = %d", got)
	}
}

func TestComparisons(t *testing.T) {
	usd := Currency{Code: "USD", Scale: 2, Symbol: "$"}
	a := New(500, usd)
	b := New(700, usd)

	gt, err := b.GreaterThan(a)
	if err != nil {
		t.Fatalf("greater than error: %v", err)
	}
	if !gt {
		t.Fatalf("expected b > a")
	}
	lt, err := a.LessThan(b)
	if err != nil {
		t.Fatalf("less than error: %v", err)
	}
	if !lt {
		t.Fatalf("expected a < b")
	}
	if !a.Equal(New(500, usd)) {
		t.Fatalf("expected equal")
	}
}

func TestMul(t *testing.T) {
	usd := Currency{Code: "USD", Scale: 2, Symbol: "$"}
	m := New(1050, usd)

	out, err := m.Mul(2)
	if err != nil {
		t.Fatalf("mul error: %v", err)
	}
	if got := out.Amount(); got != 2100 {
		t.Fatalf("mul amount = %d", got)
	}
}

func TestDiv(t *testing.T) {
	usd := Currency{Code: "USD", Scale: 2, Symbol: "$"}
	m := New(2100, usd)

	out, err := m.Div(3)
	if err != nil {
		t.Fatalf("div error: %v", err)
	}
	if got := out.Amount(); got != 700 {
		t.Fatalf("div amount = %d", got)
	}
}

func TestStringFormatting(t *testing.T) {
	usd := Currency{Code: "USD", Scale: 2, Symbol: "$"}
	m := New(-105, usd)
	if got := m.String(); got != "-$1.05" {
		t.Fatalf("string = %s", got)
	}
	jpy := Currency{Code: "JPY", Scale: 0, Symbol: "¥"}
	m = New(123, jpy)
	if got := m.String(); got != "¥123" {
		t.Fatalf("string = %s", got)
	}
}

func TestFormatConfig(t *testing.T) {
	orig := DefaultFormat()
	defer func() {
		if err := SetFormat(orig); err != nil {
			t.Fatalf("reset format: %v", err)
		}
	}()

	cfg := FormatConfig{
		DecimalSeparator:   ",",
		ThousandsSeparator: ".",
		SymbolPosition:     SymbolSuffix,
		SymbolKind:         SymbolUseCurrencySymbol,
		Space:              true,
	}

	if err := SetFormat(cfg); err != nil {
		t.Fatalf("set format: %v", err)
	}

	try := Currency{Code: "TRY", Scale: 2, Symbol: "₺"}
	m := New(21229, try)
	if got := m.String(); got != "212,29 ₺" {
		t.Fatalf("string = %s", got)
	}

	alt := FormatConfig{
		DecimalSeparator:   ",",
		ThousandsSeparator: " ",
		SymbolPosition:     SymbolSuffix,
		SymbolKind:         SymbolUseCurrencyCode,
		Space:              true,
	}
	text, err := m.Format(alt)
	if err != nil {
		t.Fatalf("format: %v", err)
	}
	if text != "212,29 TRY" {
		t.Fatalf("format = %s", text)
	}
}

func TestFormatGrouping(t *testing.T) {
	cfg := FormatConfig{
		DecimalSeparator:   ",",
		ThousandsSeparator: ".",
		SymbolPosition:     SymbolSuffix,
		SymbolKind:         SymbolUseCurrencySymbol,
		Space:              true,
	}

	eur := Currency{Code: "EUR", Scale: 2, Symbol: "€"}
	m := New(123456789, eur)
	text, err := m.Format(cfg)
	if err != nil {
		t.Fatalf("format: %v", err)
	}
	if text != "1.234.567,89 €" {
		t.Fatalf("format = %s", text)
	}
}
