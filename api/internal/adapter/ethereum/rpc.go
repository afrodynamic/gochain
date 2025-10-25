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

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
	ID     int             `json:"id"`
}

func (c *rpcClient) call(ctx context.Context, method string, params []any, out interface{}) error {
	if c.url == "" {
		return fmt.Errorf("ETH_RPC not set")
	}

	if params == nil {
		params = []any{}
	}

	reqBody, err := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(reqBody))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rpc %s: http %d", method, resp.StatusCode)
	}

	var res rpcResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	if res.Error != nil {
		return fmt.Errorf("rpc %s: (%d) %s", method, res.Error.Code, res.Error.Message)
	}

	if out == nil {
		return nil
	}

	return json.Unmarshal(res.Result, out)
}
