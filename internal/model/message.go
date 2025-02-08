package model

import (
	"time"

	"github.com/google/uuid"
)

type messageKey struct {
	ID         uuid.UUID
	MessageID  uuid.UUID
	SecretPart []byte
	UpdatedAt  time.Time
}

type Message struct {
	ID               uuid.UUID
	CreatorName      string
	EncryptedContent []byte
	KeyHash          []byte
	CreatedAt        time.Time
	Keys             []messageKey
}

func NewMessage(nickname string, encryptedMessage []byte, keyHash [32]byte, keyholdersInt int, keyShares [][]byte) *Message {
	message := Message{
		ID:               uuid.New(),
		CreatorName:      nickname,
		EncryptedContent: encryptedMessage,
		KeyHash:          keyHash[:],
		CreatedAt:        time.Now(),
	}

	messageKeys := make([]messageKey, keyholdersInt)
	for i := 0; i < keyholdersInt; i++ {
		messageKeys[i] = messageKey{
			ID:         uuid.New(),
			MessageID:  message.ID,
			SecretPart: keyShares[i],
			UpdatedAt:  message.CreatedAt,
		}
	}
	message.Keys = messageKeys

	return &message
}
