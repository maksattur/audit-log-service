package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"time"
)

// Config is configuration for audit log service
type Config struct {
	HTTPAddr        string        `envconfig:"HTTP_ADDR" required:"true" default:":8080"`
	KafkaBrokerAddr string        `envconfig:"KAFKA_BROKER_ADDR" required:"true" default:"kafka:9092"`
	KafkaGroupID    string        `envconfig:"KAFKA_GROUP_ID" required:"true" default:"my-group"`
	KafkaTopic      string        `envconfig:"KAFKA_TOPIC" required:"true" default:"event-log"`
	SecretKey       []byte        `envconfig:"SECRET_KEY" required:"true" default:"jksdiweJask"`
	JwtTTL          time.Duration `envconfig:"JWT_TTL" default:"100h" required:"true"`
	Postgres        Postgres
	ClickHouse      ClickHouse
}

type Postgres struct {
	DSN string `envconfig:"POSTGRES_DSN" default:"postgres://postgres:qwerty123@localhost:5432/audit?connect_timeout=5&sslmode=disable" required:"true"`
}

type ClickHouse struct {
	ClickHouseAddr string `envconfig:"CLICKHOUSE_ADDR" required:"true" default:"localhost:19000"`
	Username       string `envconfig:"CLICKHOUSE_USERNAME" required:"true" default:"default"`
	Password       string `envconfig:"CLICKHOUSE_PASSWORD" required:"true" default:"qwerty123"`
	Database       string `envconfig:"CLICKHOUSE_DATABASE" required:"true" default:"default"`
}

func NewConfigFromEnv() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("process config from env: %w", err)
	}

	return &cfg, nil
}
