package adapter

import "context"

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
