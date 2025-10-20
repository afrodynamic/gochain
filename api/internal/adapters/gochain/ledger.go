package gochain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/afrodynamic/gochain/api/internal/core"
)

type utxo struct {
	Addr   string
	Amount uint64
}

type block struct {
	Height int
	Prev   string
	Tx     core.Tx
	Hash   string
	Time   int64
}

type state struct {
	mu     sync.RWMutex
	utxos  map[string][]utxo
	blocks []block
	nonce  map[string]uint64
}

func newState() *state {
	return &state{
		utxos:  make(map[string][]utxo),
		blocks: []block{{Height: 0, Prev: "", Hash: "genesis", Time: time.Now().Unix()}},
		nonce:  make(map[string]uint64),
	}
}

func (s *state) balance(addr string) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var sum uint64
	for _, u := range s.utxos[addr] {
		sum += u.Amount
	}
	return sum
}

func (s *state) apply(tx core.Tx) (string, error) {
	if tx.Amount == 0 {
		return "", errors.New("amount must be > 0")
	}
	if s.balance(tx.From) < tx.Amount+tx.Fee {
		return "", errors.New("insufficient funds")
	}
	if tx.Nonce != s.nonce[tx.From] {
		return "", errors.New("invalid nonce")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	remain := tx.Amount + tx.Fee
	next := make([]utxo, 0, len(s.utxos[tx.From]))
	for _, u := range s.utxos[tx.From] {
		if remain == 0 {
			next = append(next, u)
			continue
		}
		if u.Amount <= remain {
			remain -= u.Amount
		} else {
			u.Amount -= remain
			remain = 0
			next = append(next, u)
		}
	}
	s.utxos[tx.From] = next
	s.utxos[tx.To] = append(s.utxos[tx.To], utxo{Addr: tx.To, Amount: tx.Amount})
	s.nonce[tx.From]++

	raw, _ := json.Marshal(tx)
	h := sha256.Sum256(raw)
	hash := hex.EncodeToString(h[:])
	b := block{Height: len(s.blocks), Prev: s.blocks[len(s.blocks)-1].Hash, Tx: tx, Hash: hash, Time: time.Now().Unix()}
	s.blocks = append(s.blocks, b)
	return hash, nil
}
