package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config is configuration for audit log service
type Config struct {
	HTTPAddr string `envconfig:"HTTP_ADDR" required:"true" default:":8080"`
}

func NewConfigFromEnv() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("process config from env: %w", err)
	}

	return &cfg, nil
}
