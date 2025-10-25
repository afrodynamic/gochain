package ethereum

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCall_NoRPCEnv(t *testing.T) {
	t.Parallel()

	client := newRPC()
	client.url = ""

	var result any
	err := client.call(context.Background(), "web3_clientVersion", nil, &result)

	if err == nil {
		t.Fatal("expected error when ETH_RPC not set")
	}
}

func TestCall_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "TestClient/1.0",
		})
	}))

	defer server.Close()

	client := newRPC()
	client.url = server.URL

	var result string

	if err := client.call(context.Background(), "web3_clientVersion", []any{}, &result); err != nil {
		t.Fatal(err)
	}

	if result != "TestClient/1.0" {
		t.Fatalf("got=%s", result)
	}
}

func TestCall_RPCError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"error":   map[string]any{"code": -32000, "message": "boom"},
		})
	}))

	defer server.Close()

	client := newRPC()
	client.url = server.URL

	var result any

	if err := client.call(context.Background(), "eth_getBalance", []any{"0x0", "latest"}, &result); err == nil {
		t.Fatal("expected rpc error")
	}
}
