package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/everlock/internal/model"
)

func SaveNewMessage(p *pgxpool.Pool, m *model.Message) error {
	tx, err := p.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to start transaction, error: %v", err)
	}

	var txErr error

	defer func() {
		if txErr != nil {
			if rbErr := tx.Rollback(context.Background()); rbErr != nil {
				log.Printf("Failed to rollback transaction, error: %v", rbErr)
			}
		}
	}()

	_, txErr = tx.Exec(
		context.Background(),
		`INSERT INTO message (id, creator_name, encrypted_content, key_hash) VALUES ($1, $2, $3, $4)`,
		m.ID, m.CreatorName, m.EncryptedContent, m.KeyHash,
	)
	if txErr != nil {
		return fmt.Errorf("failed to save message %v, error: %v", m, txErr)
	}

	for _, key := range m.Keys {
		_, txErr = tx.Exec(
			context.Background(),
			`INSERT INTO message_key (id, message_id, secret_part, updated_at) VALUES ($1, $2, $3, $4)`,
			key.ID, key.MessageID, key.SecretPart, key.UpdatedAt,
		)
		if txErr != nil {
			return fmt.Errorf("failed to save key share %v, error: %v", key, txErr)
		}
	}

	txErr = tx.Commit(context.Background())
	if txErr != nil {
		return fmt.Errorf("failed to commit transaction, error: %v", txErr)
	}

	return nil
}
