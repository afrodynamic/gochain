package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/afrodynamic/gochain/api/internal/core"
)

type Service struct {
	defaultAdapter core.ChainAdapter
	registry       *core.Registry
}

func NewService(def core.ChainAdapter, reg *core.Registry) *Service {
	return &Service{defaultAdapter: def, registry: reg}
}

func (s *Service) Adapters() []string {
	return s.registry.List()
}

func (s *Service) AdapterFor(name string) (core.ChainAdapter, error) {
	if name == "" {
		return s.defaultAdapter, nil
	}
	a, ok := s.registry.Get(name)
	if !ok {
		return nil, errors.New("unknown chain")
	}
	return a, nil
}

func (s *Service) Balance(ctx context.Context, a core.ChainAdapter, address string) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.Balance(ctx, address)
}
