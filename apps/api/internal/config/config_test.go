package config

import "testing"

func requiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/waytutor?sslmode=disable")
	t.Setenv("BASE_URL", "http://localhost:8080/")
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	t.Setenv("OAUTH_STATE_SECRET", "0123456789abcdef")
	t.Setenv("AI_GATEWAY_API_KEY", "gateway-key")
}

func TestLoadAppliesDefaultsAndNormalizesBaseURL(t *testing.T) {
	requiredEnv(t)
	t.Setenv("ENV", "")
	t.Setenv("ADDR", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("MONTHLY_CHARS_PER_ACCOUNT", "")
	t.Setenv("GLOBAL_MONTHLY_CHARS", "")
	t.Setenv("RATE_PER_MINUTE", "")

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
	if cfg.MonthlyCharsPerAccount != 500_000 {
		t.Errorf("MonthlyCharsPerAccount = %d, want 500000", cfg.MonthlyCharsPerAccount)
	}
	if cfg.GlobalMonthlyChars != 6_000_000 {
		t.Errorf("GlobalMonthlyChars = %d, want 6000000", cfg.GlobalMonthlyChars)
	}
	if cfg.RatePerMinute != 40 {
		t.Errorf("RatePerMinute = %d, want 40", cfg.RatePerMinute)
	}
	if cfg.R2Enabled() {
		t.Error("R2Enabled() = true, want the disk cache by default")
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
		{"gateway api key", "AI_GATEWAY_API_KEY"},
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

func TestLoadRejectsGlobalQuotaBelowTheAccountQuota(t *testing.T) {
	requiredEnv(t)
	t.Setenv("MONTHLY_CHARS_PER_ACCOUNT", "900000")
	t.Setenv("GLOBAL_MONTHLY_CHARS", "100000")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error when the global quota is below the account quota")
	}
}

func TestLoadAcceptsACompleteR2Group(t *testing.T) {
	requiredEnv(t)
	t.Setenv("R2_ENDPOINT", "https://account.r2.cloudflarestorage.com")
	t.Setenv("R2_BUCKET", "waytutor-audio")
	t.Setenv("R2_ACCESS_KEY_ID", "key")
	t.Setenv("R2_SECRET_ACCESS_KEY", "secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.R2Enabled() {
		t.Error("R2Enabled() = false, want true")
	}
}

func TestLoadRejectsAPartialR2Group(t *testing.T) {
	requiredEnv(t)
	t.Setenv("R2_BUCKET", "waytutor-audio")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for a partial R2 configuration")
	}
}

func TestLoadRejectsUnknownEnv(t *testing.T) {
	requiredEnv(t)
	t.Setenv("ENV", "staging")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for an unknown ENV")
	}
}

func TestLoadDatabaseURLRequiresValue(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if _, err := LoadDatabaseURL(); err == nil {
		t.Fatal("LoadDatabaseURL() error = nil, want an error for a missing DATABASE_URL")
	}
}
