package gochain

import (
	"errors"
	"sync"

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

func (chain *Chain) SubmitTx(_ core.Tx) ([]byte, error) {
	chain.mutex.Lock()
	previousBlock := chain.store.Blocks[len(chain.store.Blocks)-1]

	newBlock := core.Block{
		Hash:     []byte{byte(len(chain.store.Blocks))},
		Height:   uint64(len(chain.store.Blocks)),
		PrevHash: previousBlock.Hash,
	}
	chain.mutex.Unlock()

	sealedBlock, err := chain.engine.Seal(newBlock)

	if err != nil {
		return nil, err
	}

	chain.mutex.Lock()
	chain.store.Blocks = append(chain.store.Blocks, sealedBlock)
	index := len(chain.store.Blocks) - 1
	chain.mutex.Unlock()

	return []byte{byte(index)}, nil
}

func (chain *Chain) GetBalance(address []byte) (uint64, error) {
	chain.mutex.RLock()
	defer chain.mutex.RUnlock()

	return chain.store.Balances[string(address)], nil
}
