package calc

import (
	"fmt"

	"github.com/govalues/decimal"
)

type amount struct {
	dec decimal.Decimal
}

// Add returns the sum of two minor-unit amounts using the given scale.
// Example: Add(1050, 250, 2) -> 1300.
func Add(a, b int64, scale int32) (int64, error) {
	da, err := newAmount(a, scale)
	if err != nil {
		return 0, err
	}
	db, err := newAmount(b, scale)
	if err != nil {
		return 0, err
	}
	sum, err := da.add(db)
	if err != nil {
		return 0, err
	}
	return Round(sum.dec, scale)
}

// Sub returns the difference of two minor-unit amounts using the given scale.
// Example: Sub(1050, 250, 2) -> 800.
func Sub(a, b int64, scale int32) (int64, error) {
	da, err := newAmount(a, scale)
	if err != nil {
		return 0, err
	}
	db, err := newAmount(b, scale)
	if err != nil {
		return 0, err
	}
	diff, err := da.sub(db)
	if err != nil {
		return 0, err
	}
	return Round(diff.dec, scale)
}

// AddPercent applies an integer percent increase to a minor-unit amount.
// Example: AddPercent(10000, 10, 2) -> 11000.
func AddPercent(value, percent int64, scale int32) (int64, error) {
	da, err := newAmount(value, scale)
	if err != nil {
		return 0, err
	}
	out, err := da.addPercent(percent)
	if err != nil {
		return 0, err
	}
	return Round(out.dec, scale)
}

// SubtractPercent applies an integer percent decrease to a minor-unit amount.
// Example: SubtractPercent(10000, 10, 2) -> 9000.
func SubtractPercent(value, percent int64, scale int32) (int64, error) {
	da, err := newAmount(value, scale)
	if err != nil {
		return 0, err
	}
	out, err := da.subtractPercent(percent)
	if err != nil {
		return 0, err
	}
	return Round(out.dec, scale)
}

// Compare compares two minor-unit amounts using the given scale.
// Example: Compare(100, 200, 2) -> -1.
func Compare(a, b int64, scale int32) (int, error) {
	da, err := newAmount(a, scale)
	if err != nil {
		return 0, err
	}
	db, err := newAmount(b, scale)
	if err != nil {
		return 0, err
	}
	return da.dec.Cmp(db.dec), nil
}

// Mul multiplies a minor-unit amount by an integer factor.
// Example: Mul(1000, 2, 2) -> 2000.
func Mul(value, factor int64, scale int32) (int64, error) {
	da, err := newAmount(value, scale)
	if err != nil {
		return 0, err
	}
	mult, err := decimal.New(factor, 0)
	if err != nil {
		return 0, err
	}
	out, err := da.multiply(mult)
	if err != nil {
		return 0, err
	}
	return Round(out.dec, scale)
}

// Div divides a minor-unit amount by an integer divisor.
// Example: Div(1000, 2, 2) -> 500.
func Div(value, divisor int64, scale int32) (int64, error) {
	da, err := newAmount(value, scale)
	if err != nil {
		return 0, err
	}
	div, err := decimal.New(divisor, 0)
	if err != nil {
		return 0, err
	}
	out, err := da.divide(div, scale)
	if err != nil {
		return 0, err
	}
	return Round(out.dec, scale)
}

// newAmount wraps minor units into a decimal with the provided scale.
// Example: newAmount(1050, 2) -> 10.50.
func newAmount(value int64, scale int32) (amount, error) {
	d, err := decimal.New(value, int(scale))
	if err != nil {
		return amount{}, err
	}
	return amount{dec: d}, nil
}

// add returns a+b as a decimal amount.
// Example: 10.50 + 2.50 -> 13.00.
func (a amount) add(b amount) (amount, error) {
	d, err := a.dec.Add(b.dec)
	if err != nil {
		return amount{}, err
	}
	return amount{dec: d}, nil
}

// sub returns a-b as a decimal amount.
// Example: 10.50 - 2.50 -> 8.00.
func (a amount) sub(b amount) (amount, error) {
	d, err := a.dec.Sub(b.dec)
	if err != nil {
		return amount{}, err
	}
	return amount{dec: d}, nil
}

// addPercent applies a percentage increase to the amount.
// Example: 10.00 + 10% -> 11.00.
func (a amount) addPercent(percent int64) (amount, error) {
	mult, err := percentMultiplier(percent, true)
	if err != nil {
		return amount{}, err
	}
	return a.multiply(mult)
}

// subtractPercent applies a percentage decrease to the amount.
// Example: 10.00 - 10% -> 9.00.
func (a amount) subtractPercent(percent int64) (amount, error) {
	mult, err := percentMultiplier(percent, false)
	if err != nil {
		return amount{}, err
	}
	return a.multiply(mult)
}

// multiply multiplies the amount by a decimal multiplier.
// Example: 10.00 * 1.10 -> 11.00.
func (a amount) multiply(mult decimal.Decimal) (amount, error) {
	scale := a.dec.Scale() + mult.Scale()
	if scale > decimal.MaxScale {
		return amount{}, fmt.Errorf("scale overflow")
	}
	d, err := a.dec.MulExact(mult, scale)
	if err != nil {
		return amount{}, err
	}
	return amount{dec: d}, nil
}

// divide divides the amount by a decimal divisor using the target scale.
// Example: 10.00 / 2 -> 5.00.
func (a amount) divide(div decimal.Decimal, scale int32) (amount, error) {
	d, err := a.dec.QuoExact(div, int(scale))
	if err != nil {
		return amount{}, err
	}
	return amount{dec: d}, nil
}
