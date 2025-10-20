package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type rpcClient struct {
	url    string
	client *http.Client
}

func newRPC() *rpcClient {
	return &rpcClient{
		url:    os.Getenv("ETH_RPC"),
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type rpcReq struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcRes struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
	ID     int             `json:"id"`
}

func (c *rpcClient) call(ctx context.Context, method string, params []interface{}, out interface{}) error {
	if c.url == "" {
		return fmt.Errorf("ETH_RPC not set")
	}
	b, _ := json.Marshal(rpcReq{Jsonrpc: "2.0", Method: method, Params: params, ID: 1})
	req, _ := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r rpcRes
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Error != nil {
		return fmt.Errorf("rpc %s: (%d) %s", method, r.Error.Code, r.Error.Message)
	}
	if out != nil {
		return json.Unmarshal(r.Result, out)
	}
	return nil
}
