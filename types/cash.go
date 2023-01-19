package types

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	CashCent   = Cash(100)
	CashDollar = 100 * CashCent
	cashMax    = 1000000000 * int64(CashDollar) // 1 billion $
)

var (
	ErrorCashMaxed = errors.New("this tool does not support inputs larger than 1 billion dollars")
)

// The Cash type represents canadian dollars.
// It's a uint64 representing 1/100 of a cent.
type Cash int64

// Parse dollar amount from string, and up to 4 decimals.
// Optionally allows for negative values.
// Valid strings examples are:
// "-12.5$", "50", "1'000CAD", "$1000", "-12`345.2ca$"
func NewCash(cashStr string, allowNegative bool) (Cash, error) {
	i, err := parseNumber(cashStr, allowNegative, 4)
	if cashMax < -cashMax || cashMax < i {
		return 0, ErrorCashMaxed
	}
	return Cash(i), err
}

// Returns the string value of cash formatted as:
// "-1'234.56" or "1'234" if no cents present
func (c Cash) String() string {
	// Prepare the number
	isNegative := c < 0
	if isNegative {
		c *= -1
	}
	if c%100 >= 50 {
		c = c + 100
	}
	c = c / 100
	// Format cents
	cents := strconv.FormatInt(int64(c%CashCent), 10)
	if cents == "0" {
		cents = ""
	} else if len(cents) == 1 {
		cents = ".0" + cents
	} else {
		cents = "." + cents
	}
	// Format dollars
	dollarsRaw := strconv.FormatInt(int64(c/CashCent), 10)
	lenDollarsRaw := len(dollarsRaw)
	dollarsBuilder := strings.Builder{}
	dollarsBuilder.Grow(lenDollarsRaw + lenDollarsRaw/3 + 1)
	if isNegative {
		dollarsBuilder.WriteByte('-')
	}
	for i := 0; i < lenDollarsRaw; i++ {
		dollarsBuilder.WriteByte(dollarsRaw[i])
		if (lenDollarsRaw-i)%3 == 1 && i != lenDollarsRaw-1 {
			dollarsBuilder.WriteByte('\'')
		}
	}
	return dollarsBuilder.String() + cents
}

// Return a given percentage of this money
func (c Cash) Percentage(p Percentage) (Cash, error) {
	// Quick special case for 0
	if c == 0 || p == 0 {
		return 0, nil
	}
	// Detect potential overflows
	if math.MaxInt64/i64Abs(int64(p)) < i64Abs(int64(c)) {
		return 0, ErrorOverflow
	}
	return Cash(int64(c) * int64(p) / int64(percentage100)), nil
}

// Return the fraction of this Cash over another in Percentage
// If other == 0 will return 0%
func (c Cash) FractionOf(other Cash) (Percentage, error) {
	if other == 0 {
		return 0, nil
	}
	// Detect potential overflows
	if math.MaxInt64/i64Abs(int64(percentage100)) < i64Abs(int64(c)) {
		return 0, ErrorOverflow
	}
	return Percentage(c) * percentage100 / Percentage(other), nil
}
