package ethereum

import (
	"context"
	"testing"

	"github.com/afrodynamic/gochain/api/internal/core"
)

func TestBuildAndOfflineSign(t *testing.T) {
	a := NewAdapter()
	tx, err := a.BuildTx(context.Background(), "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 7, core.FeeHint{MaxFeePerGas: 1})
	if err != nil || tx.Amount != 7 {
		t.Fatalf("build err=%v tx=%+v", err, tx)
	}
	stx, err := a.SignTx("0x"+"11"+"22"+"33"+"44"+"55"+"66"+"77"+"88"+"99"+"aa"+"bb"+"cc"+"dd"+"ee"+"ff"+"00112233445566778899aabbccddeeff", tx)
	if err != nil || stx.RawHex == "" || stx.TxID == "" {
		t.Fatalf("sign err=%v stx=%+v", err, stx)
	}
}

func TestBalanceWithoutRPC(t *testing.T) {
	a := NewAdapter()
	bal, err := a.Balance(context.Background(), "0x0000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if bal != 0 {
		t.Fatalf("bal=%d", bal)
	}
}

func TestParseAddress(t *testing.T) {
	a := NewAdapter()
	if _, err := a.ParseAddress("0x123"); err == nil {
		t.Fatal("expected error")
	}
}
