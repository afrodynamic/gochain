package consensus

import "github.com/afrodynamic/gochain/api/internal/core"

type Engine interface {
	Seal(block core.Block) (core.Block, error)
	Validate(block core.Block) error
	Name() string
}
