package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afrodynamic/gochain/api/internal/adapters/gochain"
	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/afrodynamic/gochain/api/internal/http/handlers"
	"github.com/afrodynamic/gochain/api/internal/http/openapi"
	"github.com/afrodynamic/gochain/api/internal/service/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHandlerWithAdapter() (*handlers.Handler, *gochain.Adapter) {
	ad := gochain.NewAdapter()
	reg := core.NewRegistry(ad)
	svc := wallet.NewService(ad, reg)
	return handlers.NewHandler(svc), ad
}

func TestGetV1Health(t *testing.T) {
	h, _ := newTestHandlerWithAdapter()
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rec := httptest.NewRecorder()

	h.GetV1Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestGetV1Chains(t *testing.T) {
	h, _ := newTestHandlerWithAdapter()
	req := httptest.NewRequest(http.MethodGet, "/v1/chains", nil)
	rec := httptest.NewRecorder()

	h.GetV1Chains(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var chains []string
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &chains))
	assert.Contains(t, chains, "gochain")
}

func TestPostV1KeysNew(t *testing.T) {
	h, _ := newTestHandlerWithAdapter()

	body := bytes.NewBufferString(`{"seed":"abc123"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/keys/new?chain=gochain", body)
	rec := httptest.NewRecorder()

	cv := openapi.PostV1KeysNewParamsChain("gochain")
	params := openapi.PostV1KeysNewParams{Chain: &cv}

	h.PostV1KeysNew(rec, req, params)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]string
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["address"])
	assert.NotEmpty(t, resp["privateKey"])
	assert.NotEmpty(t, resp["publicKey"])
}

func TestPostV1TxBuild(t *testing.T) {
	h, _ := newTestHandlerWithAdapter()

	body := bytes.NewBufferString(`{"from":"A","to":"B","amount":10,"fee":{"maxFeePerGas":1}}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/tx/build?chain=gochain", body)
	rec := httptest.NewRecorder()

	cv := openapi.PostV1TxBuildParamsChain("gochain")
	params := openapi.PostV1TxBuildParams{Chain: &cv}

	h.PostV1TxBuild(rec, req, params)

	assert.Equal(t, http.StatusOK, rec.Code)
	var tx core.Tx
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &tx))
	assert.Equal(t, "A", tx.From)
	assert.Equal(t, "B", tx.To)
	assert.EqualValues(t, 10, tx.Amount)
}

func TestGetV1BalanceAddress(t *testing.T) {
	h, ad := newTestHandlerWithAdapter()

	// create a funded address using the same adapter state
	_, _, addr, err := ad.NewKey([]byte("alice"))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/v1/balance/"+addr+"?chain=gochain", nil)
	rec := httptest.NewRecorder()

	cv := openapi.GetV1BalanceAddressParamsChain("gochain")
	params := openapi.GetV1BalanceAddressParams{Chain: &cv}

	h.GetV1BalanceAddress(rec, req, addr, params)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]any
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Contains(t, resp, "balance")
}
