package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/afrodynamic/gochain/api/internal/core"
)

type blockPayload struct {
	Hash         string               `json:"hash"`
	Height       uint64               `json:"height"`
	PrevHash     string               `json:"prevHash"`
	Timestamp    time.Time            `json:"timestamp"`
	Transactions []transactionPayload `json:"transactions"`
}

type transactionPayload struct {
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Amount      uint64    `json:"amount"`
	Fee         uint64    `json:"fee"`
	Nonce       uint64    `json:"nonce"`
	BlockHeight uint64    `json:"blockHeight"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
}

func withDemoCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")

		if origin != "" {
			responseWriter.Header().Set("Access-Control-Allow-Origin", origin)
			responseWriter.Header().Set("Vary", "Origin")
			responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
			responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		}

		if request.Method == http.MethodOptions {
			responseWriter.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}

func newBlocksHandler(blockchain core.Blockchain) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		limit := parseLimit(request, 50)

		blocks, err := blockchain.ListBlocks(uint64(limit))

		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

			return
		}

		payload := make([]blockPayload, 0, len(blocks))

		for _, block := range blocks {
			blockEntry := blockPayload{
				Hash:         encodeHex(block.Hash),
				Height:       block.Height,
				PrevHash:     encodeHex(block.PrevHash),
				Timestamp:    block.Timestamp,
				Transactions: make([]transactionPayload, 0, len(block.Transactions)),
			}

			for _, tx := range block.Transactions {
				blockEntry.Transactions = append(blockEntry.Transactions, convertTx(tx))
			}

			payload = append(payload, blockEntry)
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"blocks": payload,
		})
	})
}

func newTransactionsHandler(blockchain core.Blockchain) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		limit := parseLimit(request, 100)

		transactions, err := blockchain.ListTransactions(uint64(limit))

		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

			return
		}

		payload := make([]transactionPayload, 0, len(transactions))

		for _, tx := range transactions {
			payload = append(payload, convertTx(tx))
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"transactions": payload,
		})
	})
}

func parseLimit(request *http.Request, fallback int) int {
	limitParam := request.URL.Query().Get("limit")

	if limitParam == "" {
		return fallback
	}

	if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 {
		if parsed > 500 {
			return 500
		}

		return parsed
	}

	return fallback
}

func convertTx(tx core.Transaction) transactionPayload {
	return transactionPayload{
		Hash:        encodeHex(tx.Hash),
		From:        formatAddress(tx.From),
		To:          formatAddress(tx.To),
		Amount:      tx.Amount,
		Fee:         tx.Fee,
		Nonce:       tx.Nonce,
		BlockHeight: tx.BlockHeight,
		Timestamp:   tx.Timestamp,
		Status:      string(tx.Status),
	}
}

func encodeHex(value []byte) string {
	if len(value) == 0 {
		return ""
	}

	return "0x" + hex.EncodeToString(value)
}

func formatAddress(value []byte) string {
	if len(value) == 0 {
		return ""
	}

	printable := true

	for _, b := range value {
		if b < 0x20 || b > 0x7e {
			printable = false

			break
		}
	}

	if printable {
		return string(value)
	}

	return encodeHex(value)
}
