package tools

import (
	"crypto/sha256"
	"fmt"
)

func GetTokenHash(token string) ([]byte, error) {
	h := sha256.New()
	tokenBytes := []byte(token)

	if _, err := h.Write(tokenBytes); err != nil {
		return []byte{}, fmt.Errorf("не удалось получить хэш refresh токена: %w", err)
	}

	return h.Sum(nil), nil
}
