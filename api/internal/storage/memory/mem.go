package memory

import "github.com/afrodynamic/gochain/api/internal/core"

type Store struct {
	Blocks   []core.Block
	Balances map[string]uint64
}

func New() *Store {
	return &Store{
		Blocks: []core.Block{
			{
				Hash:   []byte("genesis"),
				Height: 0,
			},
		},
		Balances: make(map[string]uint64),
	}
}
