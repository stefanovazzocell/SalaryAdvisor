package types_test

import (
	"testing"

	"github.com/stefanovazzocell/SalaryAdvisor/types"
)

func TestPercentage(t *testing.T) {
	testCases := map[struct {
		str           string
		allowNegative bool
	}]struct {
		value     int32  // The parsed value we expect
		err       bool   // Should we expect an error? If so, do not test strFormat/reverse
		strFormat string // String Value
		reverse   bool   // Should we attempt a Percentage -> String -> Percentage conversion?
	}{
		{"", true}:          {0, false, "0", true},
		{"0", true}:         {0, false, "0", true},
		{".", true}:         {0, false, "0", true},
		{".0", true}:        {0, false, "0", true},
		{"-.1", true}:       {-100000, false, "-0.1", true},
		{"-1%", false}:      {1000000, false, "1", true},
		{"123%", true}:      {123000000, false, "123", true},
		{"100%", true}:      {100000000, false, "100", true},
		{"-100", true}:      {-100000000, false, "-100", true},
		{"12.345678", true}: {12345678, false, "12.35", false},
		{"0.01", true}:      {10000, false, "0.01", true},
	}

	for testParams, testResults := range testCases {
		percentage, err := types.NewPercentage(testParams.str, testParams.allowNegative)
		// Check returned value
		if percentage != types.Percentage(testResults.value) {
			t.Fatalf("Expected %d for %v, instead got %d", testResults.value, testParams, percentage)
		}
		// Check for errors
		if testResults.err {
			if err == nil {
				t.Fatalf("Expected an error for %v, instead got %v", testParams, percentage)
			} else {
				continue
			}
		}
		if err != nil {
			t.Fatalf("Got an unexpected error for %v: %v", testParams, err)
		}
		// Check string conversion
		str := percentage.String()
		if str != testResults.strFormat {
			t.Fatalf("Expected %q from String() of %v, instead got %q", testResults.strFormat, testParams, str)
		}
		// Run reverse test
		derivedPercentage, err := types.NewPercentage(str, testParams.allowNegative)
		if err != nil {
			t.Fatalf("Got error during reverse test for %v: %v", testParams, err)
		}
		if !testResults.reverse {
			// Skip the next part if we're not supposed to perform the reverse test
			continue
		}
		if derivedPercentage != percentage {
			t.Fatalf("Expected %d for the reverse of %v (in String(): %q), instead got %d", percentage, testParams, str, derivedPercentage)
		}
	}
}

func BenchmarkNewPercentage(b *testing.B) {
	b.Run("zero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			types.NewPercentage("0", true)
		}
	})
	b.Run("standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			types.NewPercentage("12.3%", false)
		}
	})
	b.Run("complex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			types.NewPercentage("12.345678 %", true)
		}
	})
}

func BenchmarkPercentageString(b *testing.B) {
	b.Run("zero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = types.Percentage(0).String()
		}
	})
	b.Run("standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = types.Percentage(12300000).String()
		}
	})
	b.Run("complex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = types.Percentage(12345678).String()
		}
	})
}

func FuzzPercentage(f *testing.F) {
	seeds := []string{
		"",
		"0",
		"0.0",
		"0.01",
		"19%",
		"38.54%",
		"12.336423%",
		"100%",
		"%24.34",
		"100",
		"-",
		"-0",
		"-0.0",
		"-0.01",
		"-19%",
		"-38.54%",
		"-12.336423%",
		"-100%",
		"%-24.34",
		"-100",
	}
	for _, seed := range seeds {
		f.Add(seed, true)
		f.Add(seed, false)
	}
	f.Fuzz(func(t *testing.T, percentageStr string, allowNegative bool) {
		percentage, err := types.NewPercentage(percentageStr, allowNegative)
		if err != nil {
			t.SkipNow()
		}
		// Reverse
		reversedPercentage, err := types.NewPercentage(percentage.String(), allowNegative)
		if err != nil {
			t.Fatalf("Failed to reverse percentage: %v", err)
		}
		// We normalize the percentage value now to match expected, then compare
		t.Logf("%d is %q (from %q)", percentage, percentage.String(), percentageStr)
		if percentage%10000 >= 5000 {
			percentage = percentage + 10000
		} else if percentage%10000 <= -5000 {
			percentage = percentage - 10000
		}
		percentage = percentage / 10000 * 10000
		if reversedPercentage != percentage {
			t.Fatalf("Expected reversed percentage %d to match percentage %d", reversedPercentage, percentage)
		}
	})
}
