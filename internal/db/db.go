package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/everlock/internal/config"
)

// LoadDbPool создает пул подключений к базе данных, выполняет пинг и возвращает пул,
// если соединение с БД успешно установлено. В случае ошибки завершает программу с ошибкой.
func LoadDbPool(c *config.Db) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name)

	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Failed to create a dbConfig, error: %v", err)
	}

	dbConfig.MaxConns = c.MaxConns
	dbConfig.MinConns = c.MinConns
	dbConfig.MaxConnLifetime = c.MaxConnLifetime
	dbConfig.MaxConnIdleTime = c.MaxConnIdleTime
	dbConfig.HealthCheckPeriod = c.HealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = c.ConnectTimeout

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to create DbPool, error: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		log.Fatalf("Unable to connect to database, error: %v", err)
	}

	log.Println("DbPool successfully created and Connected")
	return pool
}
