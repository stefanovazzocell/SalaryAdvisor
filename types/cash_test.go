package types_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/stefanovazzocell/SalaryAdvisor/types"
)

func TestCash(t *testing.T) {
	testCases := map[struct {
		str           string
		allowNegative bool
	}]struct {
		value     int64  // The parsed value we expect
		err       bool   // Should we expect an error? If so, do not test strFormat/reverse
		strFormat string // String Value
		reverse   bool   // Should we attempt a Cash -> String -> Cash conversion?
	}{
		{"", false}:                                   {0, false, "0", true},
		{"0", true}:                                   {0, false, "0", true},
		{"0.0", true}:                                 {0, false, "0", true},
		{"0.01", true}:                                {100, false, "0.01", true},
		{"19$", true}:                                 {190000, false, "19", true},
		{"19.0$", true}:                               {190000, false, "19", true},
		{"19.5$", true}:                               {195000, false, "19.50", true},
		{"1'234.56$", true}:                           {12345600, false, "1'234.56", true},
		{"1'234.5678$", true}:                         {12345678, false, "1'234.57", false},
		{"1'234.5628$", true}:                         {12345628, false, "1'234.56", false},
		{"1'234.5628ca$", true}:                       {12345628, false, "1'234.56", false},
		{"CA$1'234.5628", true}:                       {12345628, false, "1'234.56", false},
		{"CA$1'234.562892", true}:                     {12345628, false, "1'234.56", false},
		{"-0", true}:                                  {0, false, "0", true},
		{"-0.0", true}:                                {0, false, "0", true},
		{"-19$", true}:                                {-190000, false, "-19", true},
		{"-19.0$", true}:                              {-190000, false, "-19", true},
		{"-19.5$", true}:                              {-195000, false, "-19.50", true},
		{"-1'234.56$", true}:                          {-12345600, false, "-1'234.56", true},
		{"-1'234.5678$", true}:                        {-12345678, false, "-1'234.57", false},
		{"-1'234.5628$", true}:                        {-12345628, false, "-1'234.56", false},
		{"-1'234.5628ca$", true}:                      {-12345628, false, "-1'234.56", false},
		{"CA$-1'234.5628", true}:                      {-12345628, false, "-1'234.56", false},
		{"CA$-1'234.562892", true}:                    {-12345628, false, "-1'234.56", false},
		{"-0", false}:                                 {0, false, "0", true},
		{"-0.0", false}:                               {0, false, "0", true},
		{"-19$", false}:                               {190000, false, "19", true},
		{"-19.0$", false}:                             {190000, false, "19", true},
		{"-19.5$", false}:                             {195000, false, "19.50", true},
		{"-1'234.56$", false}:                         {12345600, false, "1'234.56", true},
		{"-1'234.5678$", false}:                       {12345678, false, "1'234.57", false},
		{"-1'234.5628$", false}:                       {12345628, false, "1'234.56", false},
		{"-1'234.5628ca$", false}:                     {12345628, false, "1'234.56", false},
		{"CA$-1'234.5628", false}:                     {12345628, false, "1'234.56", false},
		{"CA$-1'234.562892", false}:                   {12345628, false, "1'234.56", false},
		{strconv.FormatInt(math.MaxInt64, 10), false}: {0, true, "", false},
		{strconv.FormatUint(math.MaxUint64, 10), false}:                   {0, true, "", false},
		{strconv.FormatInt(math.MaxInt64, 10) + "000", false}:             {0, true, "", false},
		{strconv.FormatUint(math.MaxUint64, 10) + "000", false}:           {0, true, "", false},
		{"0." + strconv.FormatInt(math.MaxInt64, 10), false}:              {9223, false, "0.92", false},
		{"0." + strconv.FormatUint(math.MaxUint64, 10), false}:            {1844, false, "0.18", false},
		{"0." + strconv.FormatInt(math.MaxInt64, 10) + "000000", false}:   {9223, false, "0.92", false},
		{"0." + strconv.FormatUint(math.MaxUint64, 10) + "000000", false}: {1844, false, "0.18", false},
		{"-.1", true}:             {-1000, false, "-0.10", true},
		{"1'000'000'001", false}:  {0, true, "", false},
		{"-1'000'000'001", false}: {0, true, "", false},
	}
	for testParams, testResults := range testCases {
		cash, err := types.ParseCash(testParams.str, testParams.allowNegative)
		// Check returned value
		if cash != types.Cash(testResults.value) {
			t.Fatalf("Expected %d for %v, instead got %d", testResults.value, testParams, cash)
		}
		// Check for errors
		if testResults.err {
			if err == nil {
				t.Fatalf("Expected an error for %v, instead got %v", testParams, cash)
			} else {
				continue
			}
		}
		if err != nil {
			t.Fatalf("Got an unexpected error for %v: %v", testParams, err)
		}
		// Check string conversion
		str := cash.String()
		if str != testResults.strFormat {
			t.Fatalf("Expected %q from String() of %v, instead got %q", testResults.strFormat, testParams, str)
		}
		// Run reverse test
		derivedCash, err := types.ParseCash(str, testParams.allowNegative)
		if err != nil {
			t.Fatalf("Got error during reverse test for %v: %v", testParams, err)
		}
		if !testResults.reverse {
			// Skip the next part if we're not supposed to perform the reverse test
			continue
		}
		if derivedCash != cash {
			t.Fatalf("Expected %d for the reverse of %v (in String(): %q), instead got %d", cash, testParams, str, derivedCash)
		}
	}
}

func TestCashPercentage(t *testing.T) {
	testCases := map[struct {
		cash       types.Cash
		percentage types.Percentage
	}]struct {
		cash     types.Cash
		hasError bool
	}{
		{0, 0}: {0, false}, // Base case
		{types.CashDollar, 100 * types.PercentagePoint}:          {types.CashDollar, false}, // 100% of 1$
		{types.CashCent, 100 * types.PercentagePoint}:            {types.CashCent, false},   // 100% of .01$
		{types.CashDollar, types.PercentagePoint}:                {types.CashCent, false},   // 1% of 1$
		{12345, -12 * types.PercentagePoint}:                     {-1481, false},            // -12% of 1.2345$
		{types.Cash(math.MaxInt64), 100 * types.PercentagePoint}: {0, true},                 // Overflow
	}

	for testParams, testExpected := range testCases {
		result, err := testParams.cash.Percentage(testParams.percentage)
		if err != nil && !testExpected.hasError {
			t.Errorf("Unexpected error from %d ($%s) * %d (%s%%): %v",
				testParams.cash, testParams.cash.String(),
				testParams.percentage, testParams.percentage.String(),
				err)
		} else if err == nil && testExpected.hasError {
			t.Errorf("Expected an error from %d ($%s) * %d (%s%%), instead got solution %d ($%s)",
				testParams.cash, testParams.cash.String(),
				testParams.percentage, testParams.percentage.String(),
				result, result.String())
		} else if result != testExpected.cash {
			t.Errorf("Expected %d ($%s) from %d ($%s) * %d (%s%%), instead got %d ($%s)",
				testExpected.cash, testExpected.cash.String(),
				testParams.cash, testParams.cash.String(),
				testParams.percentage, testParams.percentage.String(),
				result, result.String())
		}
	}
}

func TestCashFractionOf(t *testing.T) {
	testCases := map[struct {
		over  types.Cash
		under types.Cash
	}]struct {
		percentage types.Percentage
		hasError   bool
	}{
		{0, 0}:                                        {0, false},
		{types.CashDollar, types.CashCent}:            {10000 * types.PercentagePoint, false},
		{types.CashCent, types.CashDollar}:            {types.PercentagePoint, false},
		{types.CashDollar, types.CashDollar}:          {100 * types.PercentagePoint, false},
		{types.Cash(math.MaxInt64), types.CashDollar}: {0, true},
	}

	for testParam, testExpected := range testCases {
		fraction, err := testParam.over.FractionOf(testParam.under)
		t.Logf("Comparing %d ($%s) over %d ($%s), got %d (%s%%), %v",
			testParam.over, testParam.over.String(),
			testParam.under, testParam.under.String(),
			fraction, fraction.String(),
			err)
		if testExpected.hasError && err != nil {
			continue
		} else if testExpected.hasError && err == nil {
			t.Fatal("Expected error, but got none")
		} else if err != nil {
			t.Fatal("Got unexpected error")
		}
		if fraction != testExpected.percentage {
			t.Fatalf("Expected percentage to be %d (%s%%)",
				testExpected.percentage, testExpected.percentage.String())
		}
	}
}

func BenchmarkNewCash(b *testing.B) {
	b.Run("zero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			types.ParseCash("0$", true)
		}
	})
	b.Run("standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			types.ParseCash("123'456.78", false)
		}
	})
	b.Run("complex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			types.ParseCash("ca$ -1'000'235.4959", true)
		}
	})
}

func BenchmarkCashString(b *testing.B) {
	b.Run("zero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = types.Cash(0).String()
		}
	})
	b.Run("standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = types.Cash(1234567800).String()
		}
	})
	b.Run("complex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = types.Cash(10002354959).String()
		}
	})
}

func FuzzCash(f *testing.F) {
	seeds := []string{
		"",
		"0",
		"0.0",
		"0.01",
		"19$",
		"19.0$",
		"19.5$",
		"1'234.56$",
		"1'234.5678$",
		"1'234.5628$",
		"1'234.5628ca$",
		"CA$1'234.5628",
		"CA$1'234.562892",
		"-0",
		"-0.0",
		"-19$",
		"-19.0$",
		"-19.5$",
		"-1'234.56$",
		"-1'234.5678$",
		"-1'234.5628$",
		"-1'234.5628ca$",
		"CA$-1'234.5628",
		"CA$-1'234.562892",
		"-.1",
	}
	for _, seed := range seeds {
		f.Add(seed, true)
		f.Add(seed, false)
	}
	f.Fuzz(func(t *testing.T, cashStr string, allowNegative bool) {
		cash, err := types.ParseCash(cashStr, allowNegative)
		if err != nil {
			t.SkipNow()
		}
		// Reverse
		reversedCash, err := types.ParseCash(cash.String(), allowNegative)
		if err != nil {
			t.Fatalf("Failed to reverse cash: %v", err)
		}
		// We normalize the cash value now to match expected, then compare
		if cash%100 >= 50 {
			cash = cash + 100
		} else if cash%100 <= -50 {
			cash = cash - 100
		}
		cash = cash / 100 * 100
		if reversedCash != cash {
			t.Fatalf("Expected reversed cash %d to match cash %d", reversedCash, cash)
		}
	})
}

func FuzzCashPercentage(f *testing.F) {
	seeds := []struct {
		cash       int64
		percentage int32
	}{
		{0, 0},
		{int64(types.CashDollar), int32(types.PercentagePoint)},
		{1234567 * int64(types.CashDollar), 10 * int32(types.PercentagePoint)},
		{12345678 * int64(types.CashDollar), 100 * int32(types.PercentagePoint)},
		{1234567 * int64(types.CashCent), -32 * int32(types.PercentagePoint)},
		{10 * int64(types.CashCent), -100 * int32(types.PercentagePoint)},
	}
	for _, seed := range seeds {
		f.Add(seed.cash, seed.percentage)
	}

	f.Fuzz(func(t *testing.T, cashVal int64, percentageVal int32) {
		if percentageVal < -100*int32(types.PercentagePoint) || 100*int32(types.PercentagePoint) < percentageVal {
			t.SkipNow()
		}
		cash := types.Cash(cashVal)
		percentage := types.Percentage(percentageVal)
		expected := cashVal * int64(percentageVal) / int64(100*types.PercentagePoint)
		actual, err := cash.Percentage(percentage)
		if err != nil {
			// If we catch out-of-bounds errors, it's not a problem
			t.SkipNow()
		}
		if actual != types.Cash(expected) {
			t.Errorf("Expected %d ($%s) for %d ($%s) * %d (%s%%), instead got %d ($%s)",
				expected, types.Cash(expected).String(),
				cash, cash.String(),
				percentage, types.Percentage(percentage).String(),
				actual, actual.String())
		}
	})
}
