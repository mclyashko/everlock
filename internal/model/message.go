package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageKey struct {
	ID         uuid.UUID
	MessageID  uuid.UUID
	SecretPart []byte
	UpdatedAt  time.Time
}

type Message struct {
	ID               uuid.UUID
	CreatorName      string
	EncryptedContent []byte
	KeyHash          [32]byte
	MinKeyholders    int
	CreatedAt        time.Time
	Keys             []MessageKey
}

func NewMessage(nickname string, encryptedMessage []byte, keyHash [32]byte, keyholders int, minKeyholders int) *Message {
	message := Message{
		ID:               uuid.New(),
		CreatorName:      nickname,
		EncryptedContent: encryptedMessage,
		KeyHash:          keyHash,
		MinKeyholders:    minKeyholders,
		CreatedAt:        time.Now(),
	}

	messageKeys := make([]MessageKey, keyholders)
	for i := 0; i < keyholders; i++ {
		messageKeys[i] = MessageKey{
			ID:         uuid.New(),
			MessageID:  message.ID,
			SecretPart: nil,
			UpdatedAt:  message.CreatedAt,
		}
	}
	message.Keys = messageKeys

	return &message
}
