package gochain

import (
	"encoding/hex"
	"errors"
	"strings"
)

func decodeAddress(address string) ([]byte, error) {
	trimmed := strings.TrimSpace(address)

	if trimmed == "" {
		return nil, errors.New("address is required")
	}

	if len(trimmed)%2 == 0 {
		if decoded, err := hex.DecodeString(trimmed); err == nil {
			return decoded, nil
		}
	}

	return []byte(trimmed), nil
}
