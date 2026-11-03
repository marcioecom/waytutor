package config

import "testing"

const databaseURL = "postgres://user:pass@localhost:5432/waytutor?sslmode=disable"

func TestLoadAppliesDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", databaseURL)
	t.Setenv("ENV", "")
	t.Setenv("ADDR", "")
	t.Setenv("LOG_LEVEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Env != "development" {
		t.Errorf("Env = %q, want development", cfg.Env)
	}
	if cfg.Addr != ":8080" {
		t.Errorf("Addr = %q, want :8080", cfg.Addr)
	}
	if cfg.Level() != 0 {
		t.Errorf("Level() = %v, want info", cfg.Level())
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for a missing DATABASE_URL")
	}
}

func TestLoadRejectsUnknownEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", databaseURL)
	t.Setenv("ENV", "staging")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for an unknown ENV")
	}
}
