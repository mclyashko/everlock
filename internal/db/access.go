package db

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	if txErr = saveMessage(tx, m); txErr != nil {
		return txErr
	}

	for _, key := range m.Keys {
		if txErr = saveMessageKey(tx, key); txErr != nil {
			return txErr
		}
	}

	txErr = tx.Commit(context.Background())
	if txErr != nil {
		return fmt.Errorf("failed to commit transaction, error: %v", txErr)
	}

	return nil
}

func GetMessageByID(p *pgxpool.Pool, u *uuid.UUID) (*model.Message, error) {
	tx, err := p.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction, error: %v", err)
	}

	var txErr error

	defer func() {
		if txErr != nil {
			if rbErr := tx.Rollback(context.Background()); rbErr != nil {
				log.Printf("Failed to rollback transaction, error: %v", rbErr)
			}
		}
	}()

	message, txErr := fetchMessage(tx, u)
	if txErr != nil {
		return nil, txErr
	}

	keys, txErr := fetchMessageKeys(tx, u)
	if txErr != nil {
		return nil, err
	}

	message.Keys = keys

	txErr = tx.Commit(context.Background())
	if txErr != nil {
		return nil, fmt.Errorf("failed to commit transaction, error: %v", txErr)
	}

	return message, nil
}

func UpdateKeySecret(p *pgxpool.Pool, k *model.MessageKey) error {
	query := `
		UPDATE message_key 
		SET secret_part = $1, updated_at = now() 
		WHERE id = $2
	`

	cmdTag, err := p.Exec(context.Background(), query, k.SecretPart, k.ID)
	if err != nil {
		return fmt.Errorf("failed to update key with id %v, error: %v", k.ID, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no key found with id: %v", k.ID)
	}

	return nil
}

func CleanKeysByMsgID(p *pgxpool.Pool, mid *uuid.UUID) error {
	query := `
		UPDATE message_key
		SET secret_part = NULL
		WHERE message_id = $1
	`

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
		query,
		mid,
	)
	if txErr != nil {
		return fmt.Errorf("failed to update message_key by message_id %v, error: %v", mid, txErr)
	}

	txErr = tx.Commit(context.Background())
	if txErr != nil {
		return fmt.Errorf("failed to commit transaction, error: %v", txErr)
	}

	return nil
}

func saveMessage(tx pgx.Tx, m *model.Message) error {
	_, err := tx.Exec(
		context.Background(),
		`INSERT INTO message (id, creator_name, encrypted_content, key_hash, min_keyholders, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		m.ID, m.CreatorName, m.EncryptedContent, m.KeyHash[:], m.MinKeyholders, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save message %v, error: %v", m, err)
	}

	return nil
}

func saveMessageKey(tx pgx.Tx, key model.MessageKey) error {
	_, err := tx.Exec(
		context.Background(),
		`INSERT INTO message_key (id, message_id, secret_part, updated_at) VALUES ($1, $2, $3, $4)`,
		key.ID, key.MessageID, key.SecretPart, key.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save key share %v, error: %v", key, err)
	}

	return nil
}

func fetchMessage(tx pgx.Tx, u *uuid.UUID) (*model.Message, error) {
	var message model.Message

	query := `
		SELECT id, creator_name, encrypted_content, key_hash, min_keyholders, created_at
		FROM message
		WHERE id = $1
	`

	var hash []byte

	err := tx.QueryRow(context.Background(), query, u).Scan(
		&message.ID,
		&message.CreatorName,
		&message.EncryptedContent,
		&hash,
		&message.MinKeyholders,
		&message.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message, error: %v", err)
	}

	message.KeyHash = [32]byte(hash)

	return &message, nil
}

func fetchMessageKeys(tx pgx.Tx, u *uuid.UUID) ([]model.MessageKey, error) {
	var keys []model.MessageKey

	keysQuery := `
		SELECT id, message_id, secret_part, updated_at
		FROM message_key
		WHERE message_id = $1
	`

	rows, err := tx.Query(context.Background(), keysQuery, u)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message keys, error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key model.MessageKey
		if err = rows.Scan(&key.ID, &key.MessageID, &key.SecretPart, &key.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message key, error: %v", err)
		}
		keys = append(keys, key)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return keys, nil
}
