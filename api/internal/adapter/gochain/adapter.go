package gochain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/afrodynamic/gochain/api/internal/adapter"
)

type Adapter struct {
	st *state
}

func NewAdapter() *Adapter {
	return &Adapter{st: newState()}
}

func (ad *Adapter) Network() string {
	return "gochain"
}

func (ad *Adapter) NewKey(seed []byte) (privateKey, publicKey, address string, err error) {
	if len(seed) == 0 {
		seed = GenerateRandomSeed()
	}

	privateKey, publicKey, address = NewKey(seed)

	ad.st.mu.Lock()
	ad.st.utxos[address] = append(ad.st.utxos[address], utxo{Address: address, Amount: 100})
	ad.st.mu.Unlock()

	return
}

func (ad *Adapter) ParseAddress(address string) (string, error) {
	if len(address) == 40 || len(address) == 64 {
		return address, nil
	}

	return "", errors.New("invalid gochain address")
}

func (ad *Adapter) Balance(ctx context.Context, address string) (uint64, error) {
	return ad.st.balance(address), nil
}

func (ad *Adapter) BuildTx(ctx context.Context, sender, recipient string, amount uint64, feeHint adapter.FeeHint) (adapter.Tx, error) {
	fee := feeHint.MaxFeePerGas

	if fee == 0 {
		fee = 1
	}

	return adapter.Tx{From: sender, To: recipient, Amount: amount, Fee: fee, Nonce: ad.st.nonce[sender]}, nil
}

func (ad *Adapter) SignTx(privateKey string, tx adapter.Tx) (adapter.SignedTx, error) {
	raw, _ := json.Marshal(tx)
	hash := sha256.Sum256(raw)

	return adapter.SignedTx{RawHex: hex.EncodeToString(raw), TxID: hex.EncodeToString(hash[:])}, nil
}

func (ad *Adapter) Broadcast(ctx context.Context, signedTx adapter.SignedTx) (string, error) {
	var tx adapter.Tx

	bytesHex, _ := hex.DecodeString(signedTx.RawHex)

	if err := json.Unmarshal(bytesHex, &tx); err != nil {
		return "", err
	}

	txID, err := ad.st.apply(tx)

	if err != nil {
		return "", err
	}

	return txID, nil
}

func (ad *Adapter) TxStatus(ctx context.Context, txID string) (adapter.Status, error) {
	return adapter.StatusMined, nil
}

var _ adapter.ChainAdapter = (*Adapter)(nil)
