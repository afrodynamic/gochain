package core_test

import (
	"testing"

	gc "github.com/afrodynamic/gochain/api/internal/adapters/gochain"
	"github.com/afrodynamic/gochain/api/internal/core"
)

func TestRegistry(t *testing.T) {
	r := core.NewRegistry(gc.NewAdapter())
	if len(r.List()) != 1 {
		t.Fatalf("len=%d", len(r.List()))
	}
	a, ok := r.Get("gochain")
	if !ok || a.Network() != "gochain" {
		t.Fatal("get gochain")
	}
}
