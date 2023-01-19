package types

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrorNumberOutOfBounds = errors.New("number out of bound, cannot parse it")
	ErrorOverflow          = errors.New("this operation got an overflow")
	allowedInFloatRegex    = regexp.MustCompile(`[^0-9\.]+`)
)

// Parse a period separated float number
// ignoring any non-digit, non-special (.-) chars.
// If allowNegative is set to false,
// the int64 returned will be an absolute value.
// Allowed decimals allows setting a number of decimal numbers,
// the returned int64 will be swifted (in base10) by this much to allow
// for the integer representation of such float.
// Notes:
// - If any dash is present in the string, it will be considered negative.
// - Any additional dots ('.') and any digits after them will be ignored
// Examples: ("ab-12.3456%cd.",false,2): 1234, ("ab-12.3456$cd.",true,4): -123456
func parseFloat(numberStr string, allowNegative bool, allowedDecimals uint8) (int64, error) {
	isNegative := allowNegative && strings.Contains(numberStr, "-")
	// Cleanup and split string
	cleanStr := allowedInFloatRegex.ReplaceAllString(numberStr, "")
	numberParts := strings.SplitN(cleanStr, ".", 3)
	// Parse unit
	number := int64(0)
	if len(numberParts) >= 1 && len(numberParts[0]) > 0 {
		i, err := strconv.ParseInt(numberParts[0], 10, 64)
		if err != nil {
			return 0, err
		}
		exp := int64(math.Pow10(int(allowedDecimals)))
		if i > math.MaxInt64/exp {
			return 0, ErrorNumberOutOfBounds
		}
		number = i * exp
	}
	// Parse decimal
	if len(numberParts) >= 2 && len(numberParts[1]) > 0 {
		diff := int(allowedDecimals) - len(numberParts[1])
		if diff < 0 {
			numberParts[1] = numberParts[1][:allowedDecimals]
		} else if diff > 0 {
			numberParts[1] = numberParts[1] + strings.Repeat("0", diff)
		}
		i, err := strconv.ParseInt(numberParts[1], 10, 64)
		if err != nil {
			return 0, err
		}
		number += i
	}
	// Negate if necessary
	if isNegative {
		number = -1 * number
	}
	return number, nil
}

// Returns the absolute value of int64
func i64Abs(x int64) int64 {
	if x >= 0 {
		return x
	}
	return -1 * x
}
