package ethereum

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCall_NoRPCEnv(t *testing.T) {
	c := newRPC()
	c.url = ""
	var out any
	err := c.call(context.Background(), "web3_clientVersion", nil, &out)
	if err == nil {
		t.Fatal("expected error when ETH_RPC not set")
	}
}

func TestCall_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "TestClient/1.0",
		})
	}))
	defer ts.Close()

	c := newRPC()
	c.url = ts.URL

	var out string
	if err := c.call(context.Background(), "web3_clientVersion", []interface{}{}, &out); err != nil {
		t.Fatal(err)
	}
	if out != "TestClient/1.0" {
		t.Fatalf("got=%s", out)
	}
}

func TestCall_RPCError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"error":   map[string]any{"code": -32000, "message": "boom"},
		})
	}))
	defer ts.Close()

	c := newRPC()
	c.url = ts.URL

	var out any
	if err := c.call(context.Background(), "eth_getBalance", []interface{}{"0x0", "latest"}, &out); err == nil {
		t.Fatal("expected rpc error")
	}
}
