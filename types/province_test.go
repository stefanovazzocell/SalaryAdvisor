package types_test

import (
	"testing"

	"github.com/stefanovazzocell/SalaryAdvisor/types"
)

func TestProvince(t *testing.T) {
	if province, ok := types.ParseProvince("AB"); province != 0 || !ok {
		t.Fatalf("Expected AB, got %d, %v", province, ok)
	}
	if province, ok := types.ParseProvince("bc"); province != 1 || !ok {
		t.Fatalf("Expected BC, got %d, %v", province, ok)
	}
	if province, ok := types.ParseProvince("Yt"); province != 12 || !ok {
		t.Fatalf("Expected YT, got %d, %v", province, ok)
	}
	if province, ok := types.ParseProvince("??"); province != 0 || ok {
		t.Fatalf("Expected unknown, got %d, %v", province, ok)
	}
}
