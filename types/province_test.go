package types_test

import (
	"bytes"
	"testing"

	"github.com/stefanovazzocell/SalaryAdvisor/types"
)

func TestProvince(t *testing.T) {
	if province, ok := types.ParseProvince("AB"); province != 0 || province.String() != "AB" || !ok {
		t.Fatalf("Expected AB, got %d (%q), %v", province, province.String(), ok)
	}
	if province, ok := types.ParseProvince("bc"); province != 1 || !ok {
		t.Fatalf("Expected BC, got %d, %v", province, ok)
	} else {
		jsonStr, err := province.MarshalText()
		if err != nil || !bytes.Equal(jsonStr, []byte("BC")) {
			t.Fatalf("Expected BC for MarshalText, got %q, %v", jsonStr, err)
		}
	}
	if province, ok := types.ParseProvince("Yt"); province != 12 || !ok {
		t.Fatalf("Expected YT, got %d, %v", province, ok)
	}
	if province, ok := types.ParseProvince("??"); province != 0 || ok {
		t.Fatalf("Expected unknown, got %d, %v", province, ok)
	}
	if str := types.Province(200).String(); str != "??" {
		t.Fatalf("Expected province unknown with string \"??\", got %q", str)
	}
}
