package gochain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/afrodynamic/gochain/api/internal/adapter"
	"github.com/afrodynamic/gochain/api/internal/core"
)

type Adapter struct {
	chain core.Blockchain
}

func NewAdapter(chain core.Blockchain) *Adapter {
	return &Adapter{chain: chain}
}

func (ad *Adapter) Network() string {
	return "gochain"
}

func (ad *Adapter) NewKey(seed []byte) (privateKey, publicKey, address string, err error) {
	if len(seed) == 0 {
		seed = GenerateRandomSeed()
	}

	privateKey, publicKey, address = NewKey(seed)

	if decoded, decodeErr := decodeAddress(address); decodeErr == nil {
		ad.chain.Credit(decoded, 100)
	}

	return
}

func (ad *Adapter) ParseAddress(address string) (string, error) {
	if address == "" {
		return "", errors.New("address is required")
	}

	if _, err := decodeAddress(address); err != nil {
		return "", err
	}

	return address, nil
}

func (ad *Adapter) Balance(ctx context.Context, address string) (uint64, error) {
	decoded, err := decodeAddress(address)

	if err != nil {
		return 0, err
	}

	return ad.chain.GetBalance(decoded)
}

func (ad *Adapter) BuildTx(ctx context.Context, sender, recipient string, amount uint64, feeHint adapter.FeeHint) (adapter.Tx, error) {
	fee := feeHint.MaxFeePerGas

	if fee == 0 {
		fee = 1
	}

	senderBytes, err := decodeAddress(sender)

	if err != nil {
		return adapter.Tx{}, err
	}

	return adapter.Tx{From: sender, To: recipient, Amount: amount, Fee: fee, Nonce: ad.chain.CurrentNonce(senderBytes)}, nil
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

	fromBytes, err := decodeAddress(tx.From)

	if err != nil {
		return "", err
	}

	toBytes, err := decodeAddress(tx.To)

	if err != nil {
		return "", err
	}

	submitted, err := ad.chain.SubmitTx(core.Tx{
		From:   fromBytes,
		To:     toBytes,
		Amount: tx.Amount,
		Fee:    tx.Fee,
		Data:   tx.Data,
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(submitted.Hash), nil
}

func (ad *Adapter) TxStatus(ctx context.Context, txID string) (adapter.Status, error) {
	return adapter.StatusMined, nil
}

var _ adapter.ChainAdapter = (*Adapter)(nil)
