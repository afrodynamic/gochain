package gochain

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/afrodynamic/gochain/api/internal/core"
)

func TestNewKey_DeterministicAndShapes(t *testing.T) {
	a := NewAdapter()
	priv1, pub1, addr1, err := a.NewKey([]byte("seed"))
	if err != nil {
		t.Fatal(err)
	}
	priv2, pub2, addr2, err := a.NewKey([]byte("seed"))
	if err != nil {
		t.Fatal(err)
	}
	if priv1 != priv2 || pub1 != pub2 || addr1 != addr2 {
		t.Fatal("keys must be deterministic for same seed")
	}

	// Validate private key
	privBytes, err := hex.DecodeString(priv1)
	if err != nil {
		t.Fatalf("invalid private key hex: %v", err)
	}
	if len(privBytes) != 64 {
		t.Fatalf("expected 64-byte private key, got %d", len(privBytes))
	}

	// Validate public key
	pubBytes, err := hex.DecodeString(pub1)
	if err != nil {
		t.Fatalf("invalid public key hex: %v", err)
	}
	if len(pubBytes) != 32 {
		t.Fatalf("expected 32-byte public key, got %d", len(pubBytes))
	}

	// Validate address
	addrBytes, err := hex.DecodeString(addr1)
	if err != nil {
		t.Fatalf("invalid address hex: %v", err)
	}
	if len(addrBytes) != 20 {
		t.Fatalf("expected 20-byte address, got %d", len(addrBytes))
	}
}

func TestEndToEnd(t *testing.T) {
	a := NewAdapter()

	_, _, addrFrom, err := a.NewKey([]byte("alice"))
	if err != nil {
		t.Fatal(err)
	}
	_, _, addrTo, err := a.NewKey([]byte("bob"))
	if err != nil {
		t.Fatal(err)
	}

	b1, err := a.Balance(context.Background(), addrFrom)
	if err != nil {
		t.Fatal(err)
	}
	if b1 == 0 {
		t.Fatal("funding")
	}

	tx, err := a.BuildTx(context.Background(), addrFrom, addrTo, 5, core.FeeHint{MaxFeePerGas: 1})
	if err != nil {
		t.Fatal(err)
	}

	stx, err := a.SignTx("priv-does-not-matter-in-sim", tx)
	if err != nil {
		t.Fatal(err)
	}
	if stx.TxID == "" || stx.RawHex == "" {
		t.Fatal("signed tx empty")
	}

	id, err := a.Broadcast(context.Background(), stx)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("txid empty")
	}

	s, err := a.TxStatus(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if s != core.StatusMined {
		t.Fatalf("status=%s", s)
	}

	bFrom, _ := a.Balance(context.Background(), addrFrom)
	bTo, _ := a.Balance(context.Background(), addrTo)
	if bTo == 0 || bFrom >= b1 {
		t.Fatal("transfer failed")
	}
}
