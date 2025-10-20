package wallet_test

import (
	"context"
	"testing"

	gc "github.com/afrodynamic/gochain/api/internal/adapters/gochain"
	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/afrodynamic/gochain/api/internal/service/wallet"
)

func TestAdapterSelection(t *testing.T) {
	reg := core.NewRegistry(gc.NewAdapter())
	svc := wallet.NewService(gc.NewAdapter(), reg)

	_, err := svc.AdapterFor("")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.AdapterFor("gochain")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.AdapterFor("unknown")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBalanceThroughService(t *testing.T) {
	a := gc.NewAdapter()
	reg := core.NewRegistry(a)
	svc := wallet.NewService(a, reg)

	priv, pub, addr, err := a.NewKey([]byte("seed"))
	if priv == "" || pub == "" || addr == "" || err != nil {
		t.Fatalf("keygen")
	}

	got, err := svc.Balance(context.Background(), a, addr)
	if err != nil {
		t.Fatal(err)
	}
	if got == 0 {
		t.Fatalf("want balance > 0, got %d", got)
	}
}
