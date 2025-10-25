package pos

import (
	"crypto/sha256"
	"errors"

	"github.com/afrodynamic/gochain/api/internal/consensus"
	"github.com/afrodynamic/gochain/api/internal/core"
)

type StakeReader interface {
	TotalStake() uint64
}

type Engine struct {
	stakeReader StakeReader
}

func New(stakeReader StakeReader) consensus.Engine {
	return &Engine{stakeReader: stakeReader}
}

func (engine *Engine) Seal(block core.Block) (core.Block, error) {
	totalStake := engine.stakeReader.TotalStake()

	if totalStake == 0 {
		return core.Block{}, errors.New("no stake")
	}

	hash := sha256.Sum256(append(block.Hash, byte(totalStake%255)))
	block.Hash = hash[:]

	return block, nil
}

func (engine *Engine) Validate(block core.Block) error {
	if len(block.Hash) == 0 {
		return errors.New("invalid proof of stake")
	}

	return nil
}

func (engine *Engine) Name() string {
	return "proof_of_stake"
}
