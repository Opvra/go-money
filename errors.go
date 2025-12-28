package money

import "errors"

var (
	// ErrCurrencyMismatch is returned when Money values use different currencies.
	// Example: New(100, USD).Add(New(100, EUR)) -> ErrCurrencyMismatch.
	ErrCurrencyMismatch = errors.New("currency mismatch")
	// ErrInvalidOperation is returned when an operation cannot be performed safely.
	// Example: overflow or invalid format configuration -> ErrInvalidOperation.
	ErrInvalidOperation = errors.New("invalid operation")
)
