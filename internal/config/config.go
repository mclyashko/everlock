package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	DevEnvironment  = "dev"
	ProdEnvironment = "prod"

	appEnvKey              = "APP_ENV"
	dbUserKey              = "DB_USER"
	dbPasswordKey          = "DB_PASSWORD"
	dbHostKey              = "DB_HOST"
	dbPortKey              = "DB_PORT"
	dbNameKey              = "DB_NAME"
	dbMaxConnsKey          = "DB_MAXCONNS"
	dbMinConnsKey          = "DB_MINCONNS"
	dbMaxConnLifetimeKey   = "DB_MAXCONNLIFETIME"
	dbMaxConnIdleTimeKey   = "DB_MAXCONNIDLETIME"
	dbHealthCheckPeriodKey = "DB_HEALTHCHECKPERIOD"
	dbConnectTimeoutKey    = "DB_CONNECTTIMEOUT"
	appPortKey             = "APP_PORT"
)

type Db struct {
	User              string
	Password          string
	Host              string
	Port              string
	Name              string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

type Web struct {
	Port string
}

type App struct {
	Db  Db
	Web Web
}

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() *App {
	env := os.Getenv(appEnvKey)

	envFile := filepath.Join("..", "..", "internal", "config", ".env.dev")
	if env == ProdEnvironment {
		envFile = filepath.Join("..", "..", "internal", "config", ".env.prod")
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading env file (%s), error: %v", envFile, err)
	}

	config := &App{
		Db: Db{
			User:              mustGetEnv(dbUserKey),
			Password:          mustGetEnv(dbPasswordKey),
			Host:              mustGetEnv(dbHostKey),
			Port:              mustGetEnv(dbPortKey),
			Name:              mustGetEnv(dbNameKey),
			MaxConns:          int32(mustGetEnvInt(dbMaxConnsKey)),
			MinConns:          int32(mustGetEnvInt(dbMinConnsKey)),
			MaxConnLifetime:   mustGetEnvDuration(dbMaxConnLifetimeKey),
			MaxConnIdleTime:   mustGetEnvDuration(dbMaxConnIdleTimeKey),
			HealthCheckPeriod: mustGetEnvDuration(dbHealthCheckPeriodKey),
			ConnectTimeout:    mustGetEnvDuration(dbConnectTimeoutKey),
		},
		Web: Web{
			Port: mustGetEnv(appPortKey),
		},
	}

	log.Println("Config successfully loaded")
	return config
}

// mustGetEnv получает строковое значение или падает
func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Fatal: missing required env variable %s", key)
	}
	return value
}

// mustGetEnvInt получает int из ENV или падает
func mustGetEnvInt(key string) int {
	valStr := mustGetEnv(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Fatal: invalid int value for %s: %s, error: %v", key, valStr, err)
	}
	return val
}

// mustGetEnvDuration получает time.Duration из ENV или падает
func mustGetEnvDuration(key string) time.Duration {
	valStr := mustGetEnv(key)
	val, err := time.ParseDuration(valStr)
	if err != nil {
		log.Fatalf("Fatal: invalid duration for %s: %s, error: %v", key, valStr, err)
	}
	return val
}
