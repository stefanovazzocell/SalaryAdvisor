package types_test

import (
	"testing"
	"time"

	"github.com/stefanovazzocell/SalaryAdvisor/types"
)

func TestYear(t *testing.T) {
	testCases := map[string]struct {
		year   types.Year
		isLeap bool
	}{
		// This first few tests are a bit weak
		"":                     {types.Year(time.Now().Year()), types.Year(time.Now().Year()).IsLeap()},
		"0":                    {types.Year(time.Now().Year()), types.Year(time.Now().Year()).IsLeap()},
		"10000000000000000000": {types.Year(time.Now().Year()), types.Year(time.Now().Year()).IsLeap()},
		"1900":                 {1900, false},
		"2020":                 {2020, true},
		"1834":                 {1834, false},
		"2012":                 {2012, true},
	}

	for testYearStr, testExpected := range testCases {
		year := types.NewYear(testYearStr)
		if year != testExpected.year {
			t.Fatalf("Expected %d for %q, instead got %d", testExpected.year, testYearStr, year)
		}
		if isLeap := year.IsLeap(); isLeap != testExpected.isLeap {
			t.Fatalf("Expected isLead %v for %q, instead got %v", testExpected.isLeap, testYearStr, isLeap)
		}
	}
}
