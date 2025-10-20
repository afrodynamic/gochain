package gochain

import (
	"encoding/hex"
	"testing"

	"github.com/afrodynamic/gochain/api/internal/core"
)

func TestStateBalanceAndApply(t *testing.T) {
	s := newState()
	addrFrom := hex.EncodeToString([]byte{1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, 32})
	addrTo := hex.EncodeToString([]byte{32, 31, 30, 29, 28, 27, 26, 25,
		24, 23, 22, 21, 20, 19, 18, 17,
		16, 15, 14, 13, 12, 11, 10, 9,
		8, 7, 6, 5, 4, 3, 2, 1})

	s.utxos[addrFrom] = []utxo{{Addr: addrFrom, Amount: 15}}
	s.nonce[addrFrom] = 0

	if b := s.balance(addrFrom); b != 15 {
		t.Fatalf("balance from=%d", b)
	}

	tx := core.Tx{From: addrFrom, To: addrTo, Amount: 10, Fee: 1, Nonce: 0}
	id, err := s.apply(tx)
	if err != nil || id == "" {
		t.Fatalf("apply err=%v id=%s", err, id)
	}

	if b := s.balance(addrFrom); b != 4 {
		t.Fatalf("from bal=%d", b)
	}
	if b := s.balance(addrTo); b != 10 {
		t.Fatalf("to bal=%d", b)
	}
	if s.nonce[addrFrom] != 1 {
		t.Fatalf("nonce=%d", s.nonce[addrFrom])
	}
	if got := s.blocks[len(s.blocks)-1].Tx; got.From != addrFrom || got.To != addrTo {
		t.Fatalf("block tx mismatch")
	}
}

func TestStateApply_Validation(t *testing.T) {
	s := newState()
	addr := hex.EncodeToString(make([]byte, 32))
	s.utxos[addr] = []utxo{{Addr: addr, Amount: 5}}
	s.nonce[addr] = 0

	if _, err := s.apply(core.Tx{From: addr, To: addr, Amount: 0, Fee: 1, Nonce: 0}); err == nil {
		t.Fatal("want error amount=0")
	}
	if _, err := s.apply(core.Tx{From: addr, To: addr, Amount: 5, Fee: 1, Nonce: 1}); err == nil {
		t.Fatal("want error bad nonce")
	}
	if _, err := s.apply(core.Tx{From: addr, To: addr, Amount: 10, Fee: 0, Nonce: 0}); err == nil {
		t.Fatal("want error insufficient")
	}
}
