package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env                string `env:"ENV" envDefault:"development" validate:"oneof=development production"`
	Addr               string `env:"ADDR" envDefault:":8080" validate:"required"`
	BaseURL            string `env:"BASE_URL" validate:"required,url"`
	DatabaseURL        string `env:"DATABASE_URL" validate:"required,url"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID" validate:"required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET" validate:"required"`
	OAuthStateSecret   string `env:"OAUTH_STATE_SECRET" validate:"required,min=16"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}
	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	cfg.BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")
	return cfg, nil
}

// LoadDatabaseURL is what the migrate command needs: migrations touch the
// database and nothing else.
func LoadDatabaseURL() (string, error) {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return "", errors.New("DATABASE_URL is required")
	}
	return dsn, nil
}

func (c Config) Level() slog.Level {
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
