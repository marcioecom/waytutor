package config

import "testing"

func requiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/waytutor?sslmode=disable")
	t.Setenv("BASE_URL", "http://localhost:8080/")
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	t.Setenv("OAUTH_STATE_SECRET", "0123456789abcdef")
}

func TestLoadAppliesDefaultsAndNormalizesBaseURL(t *testing.T) {
	requiredEnv(t)
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
	if cfg.BaseURL != "http://localhost:8080" {
		t.Errorf("BaseURL = %q, want the trailing slash removed", cfg.BaseURL)
	}
	if cfg.Level() != 0 {
		t.Errorf("Level() = %v, want info", cfg.Level())
	}
}

func TestLoadRequiresEverySecret(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"database url", "DATABASE_URL"},
		{"base url", "BASE_URL"},
		{"google client id", "GOOGLE_CLIENT_ID"},
		{"google client secret", "GOOGLE_CLIENT_SECRET"},
		{"oauth state secret", "OAUTH_STATE_SECRET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requiredEnv(t)
			t.Setenv(tt.key, "")

			if _, err := Load(); err == nil {
				t.Fatalf("Load() error = nil, want an error for a missing %s", tt.key)
			}
		})
	}
}

func TestLoadRejectsShortOAuthStateSecret(t *testing.T) {
	requiredEnv(t)
	t.Setenv("OAUTH_STATE_SECRET", "short")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for a short OAUTH_STATE_SECRET")
	}
}

func TestLoadDatabaseURLRequiresValue(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if _, err := LoadDatabaseURL(); err == nil {
		t.Fatal("LoadDatabaseURL() error = nil, want an error for a missing DATABASE_URL")
	}
}

func TestLoadDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/waytutor?sslmode=disable")

	dsn, err := LoadDatabaseURL()
	if err != nil {
		t.Fatalf("LoadDatabaseURL() error = %v", err)
	}
	if dsn == "" {
		t.Fatal("LoadDatabaseURL() returned an empty dsn")
	}
}

func TestLoadRejectsUnknownEnv(t *testing.T) {
	requiredEnv(t)
	t.Setenv("ENV", "staging")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for an unknown ENV")
	}
}
