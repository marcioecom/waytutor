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
	Env         string `env:"ENV" envDefault:"development" validate:"oneof=development production"`
	Addr        string `env:"ADDR" envDefault:":8080" validate:"required"`
	BaseURL     string `env:"BASE_URL" validate:"required,url"`
	DatabaseURL string `env:"DATABASE_URL" validate:"required,url"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID" validate:"required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET" validate:"required"`
	OAuthStateSecret   string `env:"OAUTH_STATE_SECRET" validate:"required,min=16"`

	GatewayAPIKey string `env:"AI_GATEWAY_API_KEY" validate:"required"`
	AdminToken    string `env:"ADMIN_TOKEN"`

	CacheDir string `env:"CACHE_DIR" envDefault:"./data/audio" validate:"required"`

	R2Endpoint        string `env:"R2_ENDPOINT" validate:"omitempty,url"`
	R2Bucket          string `env:"R2_BUCKET"`
	R2AccessKeyID     string `env:"R2_ACCESS_KEY_ID"`
	R2SecretAccessKey string `env:"R2_SECRET_ACCESS_KEY"`
	R2Region          string `env:"R2_REGION" envDefault:"auto"`

	MonthlyCharsPerAccount int64 `env:"MONTHLY_CHARS_PER_ACCOUNT" envDefault:"500000" validate:"min=1"`
	GlobalMonthlyChars     int64 `env:"GLOBAL_MONTHLY_CHARS" envDefault:"6000000" validate:"gtefield=MonthlyCharsPerAccount"`
	MonthlySummaries       int64 `env:"MONTHLY_SUMMARIES_PER_ACCOUNT" envDefault:"200" validate:"min=1"`
	RatePerMinute          int   `env:"RATE_PER_MINUTE" envDefault:"40" validate:"min=1"`
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
	if err := cfg.validateR2(); err != nil {
		return Config{}, err
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

func (c Config) R2Enabled() bool {
	return c.R2Endpoint != ""
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

func (c Config) validateR2() error {
	set := 0
	for _, value := range []string{c.R2Endpoint, c.R2Bucket, c.R2AccessKeyID, c.R2SecretAccessKey} {
		if value != "" {
			set++
		}
	}

	if set != 0 && set != 4 {
		return errors.New("R2 needs R2_ENDPOINT, R2_BUCKET, R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY together")
	}
	return nil
}
