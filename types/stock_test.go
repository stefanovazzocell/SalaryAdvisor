package types_test

import (
	"math"
	"testing"

	"github.com/stefanovazzocell/SalaryAdvisor/types"
)

func TestStock(t *testing.T) {
	testCases := map[struct {
		value      string
		conversion string
	}]struct {
		hasError   bool // Does it return an error? If so, stop here
		valueError bool // Does .Value() return an error
		value      types.Cash
		cadValue   types.Cash
		isCad      bool
	}{
		{"", ""}:                                {false, false, 0, 0, true},
		{"1", "0"}:                              {false, false, types.CashDollar, types.CashDollar, true},
		{"1", "1"}:                              {false, false, types.CashDollar, types.CashDollar, true},
		{"100$", "1.0"}:                         {false, false, 100 * types.CashDollar, 100 * types.CashDollar, true},
		{"100$", ".5"}:                          {false, false, 100 * types.CashDollar, 50 * types.CashDollar, false},
		{"1000000000", "1000000000"}:            {false, true, 0, 0, false},
		{"10000000000", "1000000000"}:           {true, false, 0, 0, false},
		{"1000000000", "1000000000000000000"}:   {true, false, 0, 0, false},
		{"1000000000000000000", "1000000000"}:   {true, false, 0, 0, false},
		{"-", "-"}:                              {false, false, 0, 0, true},
		{"-1", "-0"}:                            {false, false, types.CashDollar, types.CashDollar, true},
		{"-1", "-1"}:                            {false, false, types.CashDollar, types.CashDollar, true},
		{"-100$", "-1.0"}:                       {false, false, 100 * types.CashDollar, 100 * types.CashDollar, true},
		{"-100$", "-.5"}:                        {false, false, 100 * types.CashDollar, 50 * types.CashDollar, false},
		{"-1000000000", "-1000000000"}:          {false, true, 0, 0, false},
		{"-10000000000", "-1000000000"}:         {true, false, 0, 0, false},
		{"-1000000000", "-1000000000000000000"}: {true, false, 0, 0, false},
		{"-1000000000000000000", "-1000000000"}: {true, false, 0, 0, false},
		{"", "-"}:                               {false, false, 0, 0, true},
		{"1", "-0"}:                             {false, false, types.CashDollar, types.CashDollar, true},
		{"1", "-1"}:                             {false, false, types.CashDollar, types.CashDollar, true},
		{"100$", "-1.0"}:                        {false, false, 100 * types.CashDollar, 100 * types.CashDollar, true},
		{"100$", "-.5"}:                         {false, false, 100 * types.CashDollar, 50 * types.CashDollar, false},
		{"1000000000", "-1000000000"}:           {false, true, 0, 0, false},
		{"10000000000", "-1000000000"}:          {true, false, 0, 0, false},
		{"1000000000", "-1000000000000000000"}:  {true, false, 0, 0, false},
		{"1000000000000000000", "-1000000000"}:  {true, false, 0, 0, false},
		{"-", ""}:                               {false, false, 0, 0, true},
		{"-1", "0"}:                             {false, false, types.CashDollar, types.CashDollar, true},
		{"-1", "1"}:                             {false, false, types.CashDollar, types.CashDollar, true},
		{"-100$", "1.0"}:                        {false, false, 100 * types.CashDollar, 100 * types.CashDollar, true},
		{"-100$", ".5"}:                         {false, false, 100 * types.CashDollar, 50 * types.CashDollar, false},
		{"-1000000000", "1000000000"}:           {false, true, 0, 0, false},
		{"-10000000000", "1000000000"}:          {true, false, 0, 0, false},
		{"-1000000000", "1000000000000000000"}:  {true, false, 0, 0, false},
		{"-1000000000000000000", "1000000000"}:  {true, false, 0, 0, false},
	}

	for testParam, testExpect := range testCases {
		s, err := types.ParseStock(testParam.value, testParam.conversion)
		t.Logf("Testing with types.NewStock(%q, %q)", testParam.value, testParam.conversion)
		if err != nil && !testExpect.hasError {
			t.Fatalf("Unexpected error: %v", err)
		} else if err == nil && testExpect.hasError {
			t.Fatalf("Expected error, but %v returned", s)
		} else if err != nil {
			continue
		}
		// Attempt value conversion
		cadVal, origVal, isCad, err := s.Value()
		if err != nil && !testExpect.valueError {
			t.Fatalf("Unexpected conversion value error: %v", err)
		} else if err == nil && testExpect.valueError {
			t.Fatalf("Expected conversion value error, but %v returned", s)
		} else if err != nil {
			continue
		}
		if isCad != testExpect.isCad {
			t.Fatalf("Expected %v for isCad, instead got %v", testExpect.isCad, isCad)
		}
		if cadVal != testExpect.cadValue {
			t.Fatalf("Expected %d (%sca$), but got %d (%sca$)",
				testExpect.cadValue, testExpect.cadValue.String(),
				cadVal, cadVal.String())
		}
		if origVal != testExpect.value {
			t.Fatalf("Expected %d (%s?$), but got %d (%s?$)",
				testExpect.value, testExpect.value.String(),
				origVal, origVal.String())
		}
	}
}

func TestStockPercentageSum(t *testing.T) {
	// Stock().Percentage is, mostly, just a wrapper for Cash().Percentage
	stock, err := types.ParseStock("1000$", "2$")
	if err != nil {
		t.Errorf("Failed to setup stock: %v", err)
	}
	stock, err = stock.Percentage(50 * types.PercentagePoint)
	if err != nil {
		t.Errorf("Failed to grab percentage stock: %v", err)
	}
	if _, value, _, _ := stock.Value(); value != 500*types.CashDollar {
		t.Errorf("Detected incorrect value %d (%s)", value, value.String())
	}
	// Attempt overflow
	stock, err = types.ParseStock("1000000000$", "")
	if err != nil {
		t.Errorf("Failed to setup overflow stock: %v", err)
	}
	stock, err = stock.Percentage(types.Percentage(math.MaxInt64))
	if err == nil {
		t.Errorf("Did not overflow Stock().Percentage, instead got %v", stock)
	}
}
