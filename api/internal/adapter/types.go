package adapter

type Tx struct {
	From   string
	To     string
	Amount uint64
	Fee    uint64
	Nonce  uint64
	Data   []byte
}

type SignedTx struct {
	RawHex string
	TxID   string
}

type FeeHint struct {
	MaxFeePerGas   uint64
	MaxPriorityFee uint64
}

type Status string

const (
	StatusPending Status = "pending"
	StatusMined   Status = "mined"
)
