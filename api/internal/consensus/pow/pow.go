package pow

import (
	"crypto/sha256"
	"errors"

	"github.com/afrodynamic/gochain/api/internal/consensus"
	"github.com/afrodynamic/gochain/api/internal/core"
)

type Engine struct {
	difficulty uint8
}

func New(difficulty uint8) consensus.Engine {
	return &Engine{difficulty: difficulty}
}

func (engine *Engine) Seal(block core.Block) (core.Block, error) {
	var nonce uint64

	for {
		hash := sha256.Sum256(append(block.Hash, byte(nonce)))

		if countLeadingZeroBits(hash[:]) >= int(engine.difficulty) {
			block.Hash = hash[:]
			return block, nil
		}

		nonce++

		if nonce == 0 {
			return core.Block{}, errors.New("nonce overflow")
		}
	}
}

func (engine *Engine) Validate(block core.Block) error {
	hash := sha256.Sum256(block.Hash)

	if countLeadingZeroBits(hash[:]) < int(engine.difficulty) {
		return errors.New("invalid proof of work")
	}

	return nil
}

func (engine *Engine) Name() string {
	return "proof_of_work"
}

func countLeadingZeroBits(bytes []byte) int {
	count := 0

	for _, currentByte := range bytes {
		if currentByte == 0 {
			count += 8
			continue
		}

		for bit := 7; bit >= 0; bit-- {
			if (currentByte>>uint(bit))&1 == 0 {
				count++
			} else {
				return count
			}
		}
	}

	return count
}
