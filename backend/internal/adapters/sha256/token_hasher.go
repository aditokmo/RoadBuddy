package sha256

import (
	"crypto/sha256"
	"encoding/hex"
)

type TokenHasher struct{}

func NewTokenHasher() *TokenHasher {
	return &TokenHasher{}
}

func (h *TokenHasher) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
