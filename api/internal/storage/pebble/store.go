package pebble

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	"github.com/cockroachdb/pebble"

	"github.com/afrodynamic/gochain/api/internal/core"
)

type Store struct {
	db           *pebble.DB
	Blocks       []core.Block
	Balances     map[string]uint64
	Transactions []core.Transaction
	Nonces       map[string]uint64
}

const (
	blocksKey       = "blocks"
	balancesKey     = "balances"
	transactionsKey = "transactions"
	noncesKey       = "nonces"
)

func New(path string) (*Store, error) {
	if path == "" {
		path = "data"
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	db, err := pebble.Open(absPath, &pebble.Options{})

	if err != nil {
		return nil, err
	}

	store := &Store{
		db:           db,
		Blocks:       make([]core.Block, 0),
		Balances:     make(map[string]uint64),
		Transactions: make([]core.Transaction, 0),
		Nonces:       make(map[string]uint64),
	}

	if err := store.load(); err != nil {
		db.Close()

		return nil, err
	}

	if len(store.Blocks) == 0 {
		store.Blocks = []core.Block{{
			Hash:      []byte("genesis"),
			Height:    0,
			Timestamp: time.Now().UTC(),
		}}

		if err := store.persistAll(); err != nil {
			db.Close()

			return nil, err
		}
	}

	return store, nil
}

func (store *Store) Close() error {
	if store.db == nil {
		return nil
	}

	return store.db.Close()
}

func (store *Store) load() error {
	if err := store.loadKey(blocksKey, &store.Blocks); err != nil {
		return err
	}

	if err := store.loadKey(balancesKey, &store.Balances); err != nil {
		return err
	}

	if err := store.loadKey(transactionsKey, &store.Transactions); err != nil {
		return err
	}

	if err := store.loadKey(noncesKey, &store.Nonces); err != nil {
		return err
	}

	if store.Balances == nil {
		store.Balances = make(map[string]uint64)
	}

	if store.Nonces == nil {
		store.Nonces = make(map[string]uint64)
	}

	if store.Blocks == nil {
		store.Blocks = make([]core.Block, 0)
	}

	if store.Transactions == nil {
		store.Transactions = make([]core.Transaction, 0)
	}

	return nil
}

func (store *Store) loadKey(key string, dest interface{}) error {
	val, closer, err := store.db.Get([]byte(key))

	if errors.Is(err, pebble.ErrNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	defer closer.Close()

	return json.Unmarshal(val, dest)
}

func (store *Store) persistAll() error {
	if err := store.saveKey(blocksKey, store.Blocks); err != nil {
		return err
	}

	if err := store.saveKey(balancesKey, store.Balances); err != nil {
		return err
	}

	if err := store.saveKey(transactionsKey, store.Transactions); err != nil {
		return err
	}

	if err := store.saveKey(noncesKey, store.Nonces); err != nil {
		return err
	}

	return nil
}

func (store *Store) saveKey(key string, value interface{}) error {
	bytes, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return store.db.Set([]byte(key), bytes, pebble.Sync)
}

func (store *Store) Save() error {
	return store.persistAll()
}
