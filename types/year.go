package types

import "time"

// A type representing a calendar year
type Year int64

// Parses a string into a year
// If 0 or error, uses current year
func ParseYear(yearStr string) Year {
	i, err := parseNumber(yearStr, false, 0)
	if i == 0 || err != nil {
		thisYear := time.Now().Year()
		if thisYear < 0 {
			thisYear = 0
		}
		return Year(thisYear)
	}
	return Year(i)
}

// Returns true if the year is a leap year
func (y Year) IsLeap() bool {
	return (y%4 == 0 && y%100 != 0 || y%400 == 0)
}
