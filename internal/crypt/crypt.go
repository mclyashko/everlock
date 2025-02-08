package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// EncryptAES шифрует переданные данные с использованием AES-GCM и возвращает nonce + cipher
func EncryptAES(plainData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	cipher := aesGCM.Seal(nil, nonce, plainData, nil)

	return append(nonce, cipher...), nil
}

// DecryptAES расшифровывает данные, зашифрованные с помощью AES-GCM
func DecryptAES(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("cipher is too short; cipher size: %d, nonce size: %d", len(encryptedData), nonceSize)
	}

	nonce, cipherText := encryptedData[:nonceSize], encryptedData[nonceSize:]

	plainData, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainData, nil
}
