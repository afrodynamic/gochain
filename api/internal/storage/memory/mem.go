package memory

import (
	"time"

	"github.com/afrodynamic/gochain/api/internal/core"
)

type Store struct {
	Blocks       []core.Block
	Balances     map[string]uint64
	Transactions []core.Transaction
	Nonces       map[string]uint64
}

func New() *Store {
	return &Store{
		Blocks: []core.Block{
			{
				Hash:      []byte("genesis"),
				Height:    0,
				Timestamp: time.Now().UTC(),
			},
		},
		Balances:     make(map[string]uint64),
		Transactions: make([]core.Transaction, 0),
		Nonces:       make(map[string]uint64),
	}
}
