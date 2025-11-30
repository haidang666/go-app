package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App AppConfig `require:"true"`
	DB  DBConfig  `require:"true"`
}

type AppConfig struct {
	Port int `envconfig:"APP_PORT" default:"8080"`
}

type DBConfig struct {
	Host         string `envconfig:"DB_HOST" required:"true"`
	Port         int    `envconfig:"DB_PORT" default:"5432"`
	DatabaseName string `envconfig:"DB_NAME" required:"true"`
	Username     string `envconfig:"DB_USERNAME" required:"true"`
	Password     string `envconfig:"DB_PASSWORD" required:"true"`
}

func Load() (*Config, error) {
	godotenv.Load()

	var cfg Config

	if err := envconfig.Process("APP", &cfg.App); err != nil {
		return nil, fmt.Errorf("load APP config: %w", err)
	}
	if err := envconfig.Process("DB", &cfg.DB); err != nil {
		return nil, fmt.Errorf("load DB config: %w", err)
	}

	return &cfg, nil
}
