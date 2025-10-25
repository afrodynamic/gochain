package bitcoin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/afrodynamic/gochain/api/internal/adapter"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (ad *Adapter) Network() string {
	return "bitcoin"
}

func (ad *Adapter) NewKey(seed []byte) (privateKey, publicKey, address string, err error) {
	hash := sha256.Sum256(seed)
	keyHex := hex.EncodeToString(hash[:])

	return keyHex, keyHex, keyHex, nil
}

func (ad *Adapter) ParseAddress(address string) (string, error) {
	if _, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams); err != nil {
		return "", errors.New("invalid bitcoin address")
	}

	return address, nil
}

func (ad *Adapter) Balance(ctx context.Context, address string) (uint64, error) {
	return 0, nil
}

func (ad *Adapter) BuildTx(ctx context.Context, senderAddress, recipientAddress string, amount uint64, feeHint adapter.FeeHint) (adapter.Tx, error) {
	fee := feeHint.MaxFeePerGas

	if fee == 0 {
		fee = 1
	}

	return adapter.Tx{From: senderAddress, To: recipientAddress, Amount: amount, Fee: fee, Nonce: 0}, nil
}

func (ad *Adapter) SignTx(privateKey string, tx adapter.Tx) (adapter.SignedTx, error) {
	rawTxBytes := []byte(tx.From + tx.To)
	txHash := sha256.Sum256(rawTxBytes)

	return adapter.SignedTx{
		RawHex: hex.EncodeToString(rawTxBytes),
		TxID:   hex.EncodeToString(txHash[:]),
	}, nil
}

func (ad *Adapter) Broadcast(ctx context.Context, signedTx adapter.SignedTx) (string, error) {
	return signedTx.TxID, nil
}

func (ad *Adapter) TxStatus(ctx context.Context, txID string) (adapter.Status, error) {
	return adapter.StatusPending, nil
}

var _ adapter.ChainAdapter = (*Adapter)(nil)
