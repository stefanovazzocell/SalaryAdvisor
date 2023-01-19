package types

import (
	"strconv"
)

const (
	PercentagePoint = Percentage(1000000)
	percentage100   = 100 * PercentagePoint
)

// The Percentage type represents canadian dollars.
// It's a uint64 representing 1/1'000'000 of a percentage point.
type Percentage int64

// Parse percentage from string, and up to 2 decimals.
// Optionally allows for negative values.
// Valid strings examples are:
// "10%", "60.5", "%100", "12.34%", "-43.29"
func NewPercentage(percentageStr string, allowNegative bool) (Percentage, error) {
	i, err := parseFloat(percentageStr, allowNegative, 6)
	return Percentage(i), err
}

// Returns the string value of percentage formatted as:
// "12.3" / "12.34" or "-12" if no fraction present
func (p Percentage) String() string {
	// Prepare the number
	isNegative := p < 0
	if isNegative {
		p *= -1
	}
	if p%10000 >= 5000 {
		p = p + 10000
	}
	p = p / 10000
	// Format decimal part
	decimal := strconv.FormatInt(int64(p%100), 10)
	if decimal == "0" {
		decimal = ""
	} else if len(decimal) == 2 && decimal[1] == '0' {
		decimal = "." + decimal[:1]
	} else if len(decimal) == 1 {
		decimal = ".0" + decimal
	} else {
		decimal = "." + decimal
	}
	// Format percentage
	if isNegative {
		return "-" + strconv.FormatInt(int64(p/100), 10) + decimal
	}
	return strconv.FormatInt(int64(p/100), 10) + decimal
}
