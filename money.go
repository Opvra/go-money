package money

import (
	"math"
	"strconv"

	"github.com/Opvra/go-money/internal/calc"
)

// New constructs an immutable Money value from minor units and a currency.
// Example: New(19990, TRY).Amount() -> 19990.
func New(amount int64, currency Currency) Money {
	return Money{amount: amount, currency: currency}
}

// Zero returns a zero amount for the given currency.
// Example: Zero(USD).Amount() -> 0.
func Zero(currency Currency) Money {
	return Money{amount: 0, currency: currency}
}

// Money represents a currency-aware monetary amount in minor units.
// Example: New(1050, USD) represents $10.50.
type Money struct {
	amount   int64
	currency Currency
}

// Amount returns the amount in minor units.
// Example: New(1050, USD).Amount() -> 1050.
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency of the money.
// Example: New(1050, USD).Currency().Code -> "USD".
func (m Money) Currency() Currency {
	return m.currency
}

// Add adds two Money values of the same currency.
// Example: New(1050, USD).Add(New(250, USD)) -> 1300.
func (m Money) Add(x Money) (Money, error) {
	if !sameCurrency(m.currency, x.currency) {
		return Money{}, ErrCurrencyMismatch
	}
	amount, err := calc.Add(m.amount, x.amount, m.currency.Scale)
	if err != nil {
		return Money{}, ErrInvalidOperation
	}
	return Money{amount: amount, currency: m.currency}, nil
}

// Sub subtracts one Money value from another of the same currency.
// Example: New(1050, USD).Sub(New(250, USD)) -> 800.
func (m Money) Sub(x Money) (Money, error) {
	if !sameCurrency(m.currency, x.currency) {
		return Money{}, ErrCurrencyMismatch
	}
	amount, err := calc.Sub(m.amount, x.amount, m.currency.Scale)
	if err != nil {
		return Money{}, ErrInvalidOperation
	}
	return Money{amount: amount, currency: m.currency}, nil
}

// AddPercent increases the Money amount by an integer percentage.
// Example: New(10000, USD).AddPercent(10) -> 11000.
func (m Money) AddPercent(percent int64) (Money, error) {
	amount, err := calc.AddPercent(m.amount, percent, m.currency.Scale)
	if err != nil {
		return Money{}, ErrInvalidOperation
	}
	return Money{amount: amount, currency: m.currency}, nil
}

// SubtractPercent decreases the Money amount by an integer percentage.
// Example: New(10000, USD).SubtractPercent(10) -> 9000.
func (m Money) SubtractPercent(percent int64) (Money, error) {
	amount, err := calc.SubtractPercent(m.amount, percent, m.currency.Scale)
	if err != nil {
		return Money{}, ErrInvalidOperation
	}
	return Money{amount: amount, currency: m.currency}, nil
}

// Equal reports whether two Money values are equal and share the same currency.
// Example: New(500, USD).Equal(New(500, USD)) -> true.
func (m Money) Equal(x Money) bool {
	if !sameCurrency(m.currency, x.currency) {
		return false
	}
	cmp, err := calc.Compare(m.amount, x.amount, m.currency.Scale)
	if err != nil {
		return false
	}
	return cmp == 0
}

// GreaterThan reports whether m is greater than x, requiring matching currencies.
// Example: New(700, USD).GreaterThan(New(500, USD)) -> true.
func (m Money) GreaterThan(x Money) (bool, error) {
	if !sameCurrency(m.currency, x.currency) {
		return false, ErrCurrencyMismatch
	}
	cmp, err := calc.Compare(m.amount, x.amount, m.currency.Scale)
	if err != nil {
		return false, ErrInvalidOperation
	}
	return cmp > 0, nil
}

// LessThan reports whether m is less than x, requiring matching currencies.
// Example: New(500, USD).LessThan(New(700, USD)) -> true.
func (m Money) LessThan(x Money) (bool, error) {
	if !sameCurrency(m.currency, x.currency) {
		return false, ErrCurrencyMismatch
	}
	cmp, err := calc.Compare(m.amount, x.amount, m.currency.Scale)
	if err != nil {
		return false, ErrInvalidOperation
	}
	return cmp < 0, nil
}

// IsZero reports whether the amount is zero.
// Example: Zero(USD).IsZero() -> true.
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive reports whether the amount is positive.
// Example: New(1, USD).IsPositive() -> true.
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// IsNegative reports whether the amount is negative.
// Example: New(-1, USD).IsNegative() -> true.
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// String returns a human-readable string with the configured formatting.
// Example (default): New(1050, USD).String() -> "$10.50".
func (m Money) String() string {
	text, err := formatWithConfig(m, DefaultFormat())
	if err != nil {
		return ""
	}
	return text
}

func sameCurrency(a, b Currency) bool {
	return a.Code == b.Code && a.Scale == b.Scale && a.Symbol == b.Symbol
}

func signPrefix(amount int64) string {
	if amount < 0 {
		return "-"
	}
	return ""
}

func absInt64String(amount int64) string {
	if amount >= 0 {
		return strconv.FormatInt(amount, 10)
	}
	if amount == math.MinInt64 {
		return strconv.FormatUint(uint64(math.MaxInt64)+1, 10)
	}
	return strconv.FormatInt(-amount, 10)
}
