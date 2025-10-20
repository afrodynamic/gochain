package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/afrodynamic/gochain/api/internal/adapters/gochain"
	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/afrodynamic/gochain/api/internal/http/openapi"
	"github.com/afrodynamic/gochain/api/internal/platform/httpx"
	"github.com/afrodynamic/gochain/api/internal/service/wallet"
)

type Handler struct {
	svc *wallet.Service
}

func NewHandler(svc *wallet.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetV1Health(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) GetV1Chains(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, h.svc.Adapters())
}

func (h *Handler) PostV1KeysNew(
	w http.ResponseWriter,
	r *http.Request,
	params openapi.PostV1KeysNewParams,
) {
	var body openapi.PostV1KeysNewJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && !errors.Is(err, io.EOF) {
		httpx.BadRequest(w, err)
		return
	}

	chain := ""
	if params.Chain != nil {
		chain = string(*params.Chain)
	}
	a, err := h.svc.AdapterFor(chain)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}

	mode := openapi.PostV1KeysNewJSONBodyMode("random")
	if body.Mode != nil && *body.Mode != "" {
		mode = *body.Mode
	}

	var seed []byte
	switch {
	case body.Seed != nil && *body.Seed != "":
		seed = []byte(*body.Seed)
	case mode == "deterministic" && body.Passphrase != nil && *body.Passphrase != "":
		seed = gochain.DeriveSeedFromPassphrase(*body.Passphrase)
	default:
		seed = gochain.GenerateRandomSeed()
	}

	priv, pub, addr, err := a.NewKey(seed)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]string{
		"privateKey": priv,
		"publicKey":  pub,
		"address":    addr,
	})
}

func (h *Handler) GetV1BalanceAddress(w http.ResponseWriter, r *http.Request, address string, params openapi.GetV1BalanceAddressParams) {
	chain := ""
	if params.Chain != nil {
		chain = string(*params.Chain)
	}
	adapter, err := h.svc.AdapterFor(chain)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	parsed, err := adapter.ParseAddress(address)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	bal, err := h.svc.Balance(r.Context(), adapter, parsed)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"balance": bal})
}

func (h *Handler) PostV1TxBuild(w http.ResponseWriter, r *http.Request, params openapi.PostV1TxBuildParams) {
	chain := ""
	if params.Chain != nil {
		chain = string(*params.Chain)
	}
	adapter, err := h.svc.AdapterFor(chain)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	var in struct {
		From   string       `json:"from"`
		To     string       `json:"to"`
		Amount uint64       `json:"amount"`
		Fee    core.FeeHint `json:"fee"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.BadRequest(w, err)
		return
	}
	tx, err := adapter.BuildTx(r.Context(), in.From, in.To, in.Amount, in.Fee)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, tx)
}

func (h *Handler) PostV1TxSign(w http.ResponseWriter, r *http.Request, params openapi.PostV1TxSignParams) {
	chain := ""
	if params.Chain != nil {
		chain = string(*params.Chain)
	}
	adapter, err := h.svc.AdapterFor(chain)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	var in struct {
		PrivateKey string  `json:"privateKey"`
		Tx         core.Tx `json:"tx"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.BadRequest(w, err)
		return
	}
	signed, err := adapter.SignTx(in.PrivateKey, in.Tx)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, signed)
}

func (h *Handler) PostV1TxBroadcast(w http.ResponseWriter, r *http.Request, params openapi.PostV1TxBroadcastParams) {
	chain := ""
	if params.Chain != nil {
		chain = string(*params.Chain)
	}
	adapter, err := h.svc.AdapterFor(chain)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	var stx core.SignedTx
	if err := json.NewDecoder(r.Body).Decode(&stx); err != nil {
		httpx.BadRequest(w, err)
		return
	}
	id, err := adapter.Broadcast(r.Context(), stx)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"txId": id})
}

func (h *Handler) GetV1TxId(w http.ResponseWriter, r *http.Request, id string, params openapi.GetV1TxIdParams) {
	chain := ""
	if params.Chain != nil {
		chain = string(*params.Chain)
	}
	adapter, err := h.svc.AdapterFor(chain)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	status, err := adapter.TxStatus(r.Context(), id)
	if err != nil {
		httpx.BadRequest(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"status": string(status)})
}
