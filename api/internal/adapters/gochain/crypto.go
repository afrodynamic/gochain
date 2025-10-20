package gochain

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

func GenerateRandomSeed() []byte {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return b
}

func DeriveSeedFromPassphrase(pass string) []byte {
	salt := []byte("gochain.v1")
	return argon2.IDKey([]byte(pass), salt, 1, 64*1024, 2, 32)
}

func NewKey(seed []byte) (priv, pub, addr string) {
	h := sha256.Sum256(seed)
	k := ed25519.NewKeyFromSeed(h[:])
	pk := k.Public().(ed25519.PublicKey)
	pubHex := hex.EncodeToString(pk)
	ah := sha256.Sum256(pk)
	addrHex := hex.EncodeToString(ah[:20])
	return hex.EncodeToString(k), pubHex, addrHex
}
