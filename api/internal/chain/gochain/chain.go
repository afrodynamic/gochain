package gochain

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"sync"
	"time"

	"github.com/afrodynamic/gochain/api/internal/consensus"
	"github.com/afrodynamic/gochain/api/internal/core"
	"github.com/afrodynamic/gochain/api/internal/storage/memory"
)

type Chain struct {
	mutex  sync.RWMutex
	store  *memory.Store
	engine consensus.Engine
}

func New(engine consensus.Engine, store *memory.Store) *Chain {
	return &Chain{store: store, engine: engine}
}

func (chain *Chain) Start() error {
	return nil
}

func (chain *Chain) Stop() error {
	return nil
}

func (chain *Chain) GetBlock(height uint64) (core.Block, error) {
	chain.mutex.RLock()
	defer chain.mutex.RUnlock()

	if height >= uint64(len(chain.store.Blocks)) {
		return core.Block{}, errors.New("not found")
	}

	return chain.store.Blocks[height], nil
}

func (chain *Chain) ListBlocks(limit uint64) ([]core.Block, error) {
	chain.mutex.RLock()
	defer chain.mutex.RUnlock()

	total := uint64(len(chain.store.Blocks))

	if limit == 0 || limit > total {
		limit = total
	}

	start := total - limit
	result := make([]core.Block, 0, limit)

	for i := start; i < total; i++ {
		result = append(result, chain.store.Blocks[i])
	}

	return result, nil
}

func (chain *Chain) SubmitTx(tx core.Tx) (core.Transaction, error) {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()

	if len(chain.store.Blocks) == 0 {
		return core.Transaction{}, errors.New("chain not initialised")
	}

	if tx.Amount == 0 {
		return core.Transaction{}, errors.New("amount must be positive")
	}

	fromKey := string(tx.From)
	toKey := string(tx.To)
	totalDebit := tx.Amount + tx.Fee

	balance := chain.store.Balances[fromKey]

	if balance < totalDebit {
		return core.Transaction{}, errors.New("insufficient balance")
	}

	chain.store.Balances[fromKey] = balance - totalDebit
	chain.store.Balances[toKey] += tx.Amount

	nonce := chain.store.Nonces[fromKey]
	chain.store.Nonces[fromKey] = nonce + 1

	previousBlock := chain.store.Blocks[len(chain.store.Blocks)-1]
	height := previousBlock.Height + 1
	timestamp := time.Now().UTC()

	txHashInput := make([]byte, 0, len(tx.From)+len(tx.To)+32)
	txHashInput = append(txHashInput, tx.From...)
	txHashInput = append(txHashInput, tx.To...)

	amountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amountBytes, tx.Amount)
	txHashInput = append(txHashInput, amountBytes...)

	feeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(feeBytes, tx.Fee)
	txHashInput = append(txHashInput, feeBytes...)

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	txHashInput = append(txHashInput, nonceBytes...)

	tsBytes := []byte(timestamp.Format(time.RFC3339Nano))
	txHashInput = append(txHashInput, tsBytes...)

	txHash := sha256.Sum256(txHashInput)

	blockSeed := make([]byte, 0, len(previousBlock.Hash)+len(txHash)+len(tsBytes))
	blockSeed = append(blockSeed, previousBlock.Hash...)
	blockSeed = append(blockSeed, txHash[:]...)
	blockSeed = append(blockSeed, tsBytes...)

	newBlock := core.Block{
		Hash:         blockSeed,
		Height:       height,
		PrevHash:     previousBlock.Hash,
		Timestamp:    timestamp,
		Transactions: nil,
	}

	recordedTx := core.Transaction{
		Hash:        txHash[:],
		From:        append([]byte(nil), tx.From...),
		To:          append([]byte(nil), tx.To...),
		Amount:      tx.Amount,
		Fee:         tx.Fee,
		Nonce:       nonce,
		BlockHeight: height,
		Timestamp:   timestamp,
		Status:      core.TxStatusMined,
	}

	sealedBlock, err := chain.engine.Seal(newBlock)

	if err != nil {
		return core.Transaction{}, err
	}

	recordedTx.BlockHash = append([]byte(nil), sealedBlock.Hash...)
	sealedBlock.Transactions = []core.Transaction{recordedTx}

	chain.store.Blocks = append(chain.store.Blocks, sealedBlock)
	chain.store.Transactions = append(chain.store.Transactions, recordedTx)

	return recordedTx, nil
}

func (chain *Chain) ListTransactions(limit uint64) ([]core.Transaction, error) {
	chain.mutex.RLock()
	defer chain.mutex.RUnlock()

	total := uint64(len(chain.store.Transactions))

	if limit == 0 || limit > total {
		limit = total
	}

	start := total - limit
	result := make([]core.Transaction, 0, limit)

	for i := start; i < total; i++ {
		result = append(result, chain.store.Transactions[i])
	}

	return result, nil
}

func (chain *Chain) GetBalance(address []byte) (uint64, error) {
	chain.mutex.RLock()
	defer chain.mutex.RUnlock()

	return chain.store.Balances[string(address)], nil
}

func (chain *Chain) Credit(address []byte, amount uint64) {
	if amount == 0 {
		return
	}

	chain.mutex.Lock()
	defer chain.mutex.Unlock()

	chain.store.Balances[string(address)] += amount
}

func (chain *Chain) CurrentNonce(address []byte) uint64 {
	chain.mutex.RLock()
	defer chain.mutex.RUnlock()

	return chain.store.Nonces[string(address)]
}
