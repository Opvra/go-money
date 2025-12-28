package calc

import "github.com/govalues/decimal"

// percentMultiplier returns a decimal multiplier for percent adjustments.
// Example: percentMultiplier(10, true) -> 1.10.
func percentMultiplier(percent int64, add bool) (decimal.Decimal, error) {
	base := int64(100)
	if add {
		base += percent
	} else {
		base -= percent
	}
	return decimal.New(base, 2)
}
