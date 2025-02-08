package crypt

import (
	"crypto/rand"
	"fmt"
	"io"
)

// генерирует случайный AES-ключ заданной длины (например, 32 байта для AES-256)
func GenerateRandomAESKey(keySize int) ([]byte, error) {
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return nil, fmt.Errorf("invalid AES key size: %d; valid sizes: AES-128 (16 bytes), AES-192 (24 bytes), AES-256 (32 bytes)", keySize)
	}

	key := make([]byte, keySize)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
