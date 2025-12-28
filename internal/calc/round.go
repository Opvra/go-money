package calc

import (
	"errors"
	"math"

	"github.com/govalues/decimal"
)

var errOverflow = errors.New("overflow")

// Round converts a decimal to minor units using the target scale.
// Example: Round(decimal.New(12345, 3), 2) -> 1235.
func Round(d decimal.Decimal, scale int32) (int64, error) {
	return roundToMinor(d, scale)
}

// roundToMinor rounds a decimal to minor units using the scale.
// Example: roundToMinor(12.345, 2) -> 1235.
func roundToMinor(d decimal.Decimal, scale int32) (int64, error) {
	rounded := d.Round(int(scale))
	whole, frac, ok := rounded.Int64(int(scale))
	if !ok {
		return 0, errOverflow
	}
	return combineInt64(whole, frac, scale)
}

// combineInt64 assembles whole and frac parts into minor units.
// Example: combineInt64(12, 34, 2) -> 1234.
func combineInt64(whole, frac int64, scale int32) (int64, error) {
	if scale < 0 {
		return 0, errOverflow
	}
	if scale == 0 {
		return whole, nil
	}
	factor, ok := pow10Int64(scale)
	if !ok {
		return 0, errOverflow
	}
	prod, ok := mulInt64(whole, factor)
	if !ok {
		return 0, errOverflow
	}
	res, ok := addInt64(prod, frac)
	if !ok {
		return 0, errOverflow
	}
	return res, nil
}

// pow10Int64 returns 10^scale if it fits in int64.
// Example: pow10Int64(2) -> 100, true.
func pow10Int64(scale int32) (int64, bool) {
	if scale < 0 {
		return 0, false
	}
	var v int64 = 1
	for i := int32(0); i < scale; i++ {
		if v > math.MaxInt64/10 {
			return 0, false
		}
		v *= 10
	}
	return v, true
}

// mulInt64 multiplies two int64 values with overflow detection.
// Example: mulInt64(12, 100) -> 1200, true.
func mulInt64(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	if a == math.MinInt64 {
		if b == 1 {
			return math.MinInt64, true
		}
		if b == -1 {
			return 0, false
		}
	}
	if b == math.MinInt64 {
		if a == 1 {
			return math.MinInt64, true
		}
		if a == -1 {
			return 0, false
		}
	}
	absA := absInt64(a)
	absB := absInt64(b)
	if absA > uint64(math.MaxInt64)/absB {
		return 0, false
	}
	prod := int64(absA * absB)
	if (a < 0) != (b < 0) {
		prod = -prod
	}
	return prod, true
}

// addInt64 adds two int64 values with overflow detection.
// Example: addInt64(1000, 200) -> 1200, true.
func addInt64(a, b int64) (int64, bool) {
	if b > 0 && a > math.MaxInt64-b {
		return 0, false
	}
	if b < 0 && a < math.MinInt64-b {
		return 0, false
	}
	return a + b, true
}

// absInt64 returns the absolute value as uint64.
// Example: absInt64(-5) -> 5.
func absInt64(x int64) uint64 {
	if x >= 0 {
		return uint64(x)
	}
	if x == math.MinInt64 {
		return uint64(math.MaxInt64) + 1
	}
	return uint64(-x)
}
