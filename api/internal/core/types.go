package core

import "context"

type Tx struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
	Fee    uint64 `json:"fee"`
	Nonce  uint64 `json:"nonce"`
	Data   []byte `json:"data,omitempty"`
}

type SignedTx struct {
	RawHex string `json:"rawHex"`
	TxID   string `json:"txId"`
}

type FeeHint struct {
	MaxFeePerGas   uint64 `json:"maxFeePerGas"`
	MaxPriorityFee uint64 `json:"maxPriorityFee"`
}

type Status string

const (
	StatusPending Status = "pending"
	StatusMined   Status = "mined"
)

type ChainAdapter interface {
	Network() string
	NewKey(seed []byte) (priv, pub, addr string, err error)
	ParseAddress(s string) (string, error)
	Balance(ctx context.Context, addr string) (uint64, error)
	BuildTx(ctx context.Context, from, to string, amt uint64, feeHint FeeHint) (Tx, error)
	SignTx(priv string, tx Tx) (SignedTx, error)
	Broadcast(ctx context.Context, stx SignedTx) (string, error)
	TxStatus(ctx context.Context, id string) (Status, error)
}
