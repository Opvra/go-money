package money

// Currency defines an ISO-4217 currency and its decimal scale.
// Example: Currency{Code: "USD", Scale: 2, Symbol: "$"}.
type Currency struct {
	Code   string
	Scale  int32
	Symbol string
}
