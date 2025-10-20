package bitcoin

import (
	"context"
	"testing"

	"github.com/afrodynamic/gochain/api/internal/core"
)

func TestNewKeyAndBuildSign(t *testing.T) {
	a := NewAdapter()
	priv, pub, addr, err := a.NewKey([]byte("seed"))
	if err != nil || priv == "" || pub == "" || addr == "" {
		t.Fatalf("keygen err=%v", err)
	}
	tx, err := a.BuildTx(context.Background(), addr, addr, 1, core.FeeHint{})
	if err != nil || tx.Amount != 1 {
		t.Fatalf("build err=%v tx=%+v", err, tx)
	}
	stx, err := a.SignTx(priv, tx)
	if err != nil || stx.TxID == "" || stx.RawHex == "" {
		t.Fatalf("sign err=%v stx=%+v", err, stx)
	}
	id, err := a.Broadcast(context.Background(), stx)
	if err != nil || id == "" {
		t.Fatalf("broadcast err=%v id=%s", err, id)
	}
}
