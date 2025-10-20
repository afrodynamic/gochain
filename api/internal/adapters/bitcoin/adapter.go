package bitcoin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/btcsuite/btcd/btcutil"
)

type Adapter struct{}

func NewAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) Network() string { return "bitcoin" }

func (a *Adapter) NewKey(seed []byte) (string, string, string, error) {
	h := sha256.Sum256(seed)
	x := hex.EncodeToString(h[:])
	return x, x, x, nil
}

func (a *Adapter) ParseAddress(s string) (string, error) {
	_, err := btcutil.DecodeAddress(s, nil)
	if err != nil {
		return "", errors.New("invalid bitcoin address")
	}
	return s, nil
}

func (a *Adapter) Balance(ctx context.Context, addr string) (uint64, error) { return 0, nil }

func (a *Adapter) BuildTx(ctx context.Context, from, to string, amt uint64, feeHint core.FeeHint) (core.Tx, error) {
	fee := feeHint.MaxFeePerGas
	if fee == 0 {
		fee = 1
	}
	return core.Tx{From: from, To: to, Amount: amt, Fee: fee, Nonce: 0}, nil
}

func (a *Adapter) SignTx(priv string, tx core.Tx) (core.SignedTx, error) {
	raw := []byte(tx.From + tx.To)
	h := sha256.Sum256(raw)
	return core.SignedTx{RawHex: hex.EncodeToString(raw), TxID: hex.EncodeToString(h[:])}, nil
}

func (a *Adapter) Broadcast(ctx context.Context, stx core.SignedTx) (string, error) {
	return stx.TxID, nil
}

func (a *Adapter) TxStatus(ctx context.Context, id string) (core.Status, error) {
	return core.StatusPending, nil
}
