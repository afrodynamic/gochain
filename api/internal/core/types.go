package core

import "time"

type Block struct {
	Hash         []byte
	Height       uint64
	PrevHash     []byte
	Timestamp    time.Time
	Transactions []Transaction
}

type Tx struct {
	From   []byte
	To     []byte
	Amount uint64
	Fee    uint64
	Data   []byte
}

type TxStatus string

const (
	TxStatusPending TxStatus = "pending"
	TxStatusMined   TxStatus = "mined"
)

type Transaction struct {
	Hash        []byte
	From        []byte
	To          []byte
	Amount      uint64
	Fee         uint64
	Nonce       uint64
	BlockHash   []byte
	BlockHeight uint64
	Timestamp   time.Time
	Status      TxStatus
}

type Blockchain interface {
	Start() error
	Stop() error
	GetBlock(height uint64) (Block, error)
	ListBlocks(limit uint64) ([]Block, error)
	SubmitTx(tx Tx) (Transaction, error)
	ListTransactions(limit uint64) ([]Transaction, error)
	GetBalance(address []byte) (uint64, error)
	Credit(address []byte, amount uint64)
	CurrentNonce(address []byte) uint64
}
