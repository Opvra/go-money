package money

import "testing"

func TestPipeChain(t *testing.T) {
	try := Currency{Code: "TRY", Scale: 2, Symbol: "₺"}
	price := New(19990, try)
	shipping := New(5000, try)

	invoice, err := PipeOf(price).
		Add(shipping).
		AddPercent(20).
		Result()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	if got := invoice.Amount(); got != 29988 {
		t.Fatalf("amount = %d", got)
	}
}

func TestPipeShortCircuit(t *testing.T) {
	usd := Currency{Code: "USD", Scale: 2, Symbol: "$"}
	eur := Currency{Code: "EUR", Scale: 2, Symbol: "€"}

	_, err := PipeOf(New(100, usd)).
		Add(New(100, eur)).
		AddPercent(10).
		Result()
	if err != ErrCurrencyMismatch {
		t.Fatalf("expected ErrCurrencyMismatch, got %v", err)
	}
}
