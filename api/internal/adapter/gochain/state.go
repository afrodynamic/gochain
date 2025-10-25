package gochain

import (
	"encoding/hex"
	"sync"

	"github.com/afrodynamic/gochain/api/internal/adapter"
)

type utxo struct {
	Address string
	Amount  uint64
}

type state struct {
	mu    sync.Mutex
	utxos map[string][]utxo
	nonce map[string]uint64
}

func newState() *state {
	return &state{
		utxos: make(map[string][]utxo),
		nonce: make(map[string]uint64),
	}
}

func (st *state) balance(address string) uint64 {
	st.mu.Lock()
	defer st.mu.Unlock()

	var balance uint64

	for _, entry := range st.utxos[address] {
		balance += entry.Amount
	}

	return balance
}

func (st *state) apply(tx adapter.Tx) (string, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.nonce[tx.From]++

	txID := hex.EncodeToString([]byte(tx.From + tx.To))

	return txID, nil
}
