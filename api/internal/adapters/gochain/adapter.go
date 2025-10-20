package gochain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/afrodynamic/gochain/api/internal/core"
)

type Adapter struct{ st *state }

func NewAdapter() *Adapter { return &Adapter{st: newState()} }

func (a *Adapter) Network() string { return "gochain" }

func (a *Adapter) NewKey(seed []byte) (priv, pub, addr string, err error) {
	if len(seed) == 0 {
		seed = GenerateRandomSeed()
	}
	priv, pub, addr = NewKey(seed)
	a.st.mu.Lock()
	a.st.utxos[addr] = append(a.st.utxos[addr], utxo{Addr: addr, Amount: 100})
	a.st.mu.Unlock()
	return
}

func (a *Adapter) ParseAddress(s string) (string, error) {
	if len(s) == 40 || len(s) == 64 {
		return s, nil
	}
	return "", errors.New("invalid gochain address")
}

func (a *Adapter) Balance(ctx context.Context, addr string) (uint64, error) {
	return a.st.balance(addr), nil
}

func (a *Adapter) BuildTx(ctx context.Context, from, to string, amt uint64, feeHint core.FeeHint) (core.Tx, error) {
	fee := feeHint.MaxFeePerGas
	if fee == 0 {
		fee = 1
	}
	return core.Tx{From: from, To: to, Amount: amt, Fee: fee, Nonce: a.st.nonce[from]}, nil
}

func (a *Adapter) SignTx(priv string, tx core.Tx) (core.SignedTx, error) {
	raw, _ := json.Marshal(tx)
	h := sha256.Sum256(raw)
	return core.SignedTx{RawHex: hex.EncodeToString(raw), TxID: hex.EncodeToString(h[:])}, nil
}

func (a *Adapter) Broadcast(ctx context.Context, stx core.SignedTx) (string, error) {
	var tx core.Tx
	b, _ := hex.DecodeString(stx.RawHex)
	if err := json.Unmarshal(b, &tx); err != nil {
		return "", err
	}
	id, err := a.st.apply(tx)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (a *Adapter) TxStatus(ctx context.Context, id string) (core.Status, error) {
	return core.StatusMined, nil
}
