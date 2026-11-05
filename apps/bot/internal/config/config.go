package config

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"development" validate:"oneof=development production"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`

	DiscordToken  string   `env:"DISCORD_TOKEN" validate:"required"`
	ApplicationID string   `env:"DISCORD_APPLICATION_ID" validate:"required"`
	GuildID       string   `env:"DISCORD_GUILD_ID" validate:"required"`
	ChannelID     string   `env:"DISCORD_CHANNEL_ID" validate:"required"`
	AdminUserIDs  []string `env:"DISCORD_ADMIN_USER_IDS" envSeparator:"," validate:"min=1"`

	APIURL     string        `env:"WAYTUTOR_API_URL" envDefault:"https://waytutor-api.marcio.run" validate:"required,url"`
	AdminToken string        `env:"ADMIN_TOKEN" validate:"required"`
	PollEvery  time.Duration `env:"POLL_INTERVAL" envDefault:"5m" validate:"min=30s"`
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

	cfg.APIURL = strings.TrimSuffix(cfg.APIURL, "/")
	return cfg, nil
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
