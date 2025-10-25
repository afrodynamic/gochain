package core

type Block struct {
	Hash     []byte
	Height   uint64
	PrevHash []byte
}

type Tx struct {
	From   []byte
	To     []byte
	Amount uint64
	Data   []byte
}

type Blockchain interface {
	Start() error
	Stop() error
	GetBlock(height uint64) (Block, error)
	SubmitTx(tx Tx) ([]byte, error)
	GetBalance(address []byte) (uint64, error)
}
