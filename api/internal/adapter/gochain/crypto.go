package gochain

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

func GenerateRandomSeed() []byte {
	seed := make([]byte, 32)
	_, _ = rand.Read(seed)

	return seed
}

func DeriveSeedFromPassphrase(passphrase string) []byte {
	salt := []byte("gochain.v1")

	return argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 2, 32)
}

func NewKey(seed []byte) (privateKeyHex, publicKeyHex, addressHex string) {
	seedHash := sha256.Sum256(seed)
	privateKey := ed25519.NewKeyFromSeed(seedHash[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	publicKeyHex = hex.EncodeToString(publicKey)

	addressHash := sha256.Sum256(publicKey)
	addressHex = hex.EncodeToString(addressHash[:20])

	privateKeyHex = hex.EncodeToString(privateKey)

	return privateKeyHex, publicKeyHex, addressHex
}
