package config

import "testing"

func requiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DISCORD_TOKEN", "token")
	t.Setenv("DISCORD_APPLICATION_ID", "app")
	t.Setenv("DISCORD_GUILD_ID", "guild")
	t.Setenv("DISCORD_CHANNEL_ID", "channel")
	t.Setenv("DISCORD_ADMIN_USER_IDS", "1,2")
	t.Setenv("ADMIN_TOKEN", "admin")
}

func TestLoadAppliesDefaults(t *testing.T) {
	requiredEnv(t)
	t.Setenv("ENV", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("WAYTUTOR_API_URL", "https://waytutor-api.marcio.run/")
	t.Setenv("POLL_INTERVAL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Env != "development" {
		t.Errorf("Env = %q, want development", cfg.Env)
	}
	if cfg.Level() != 0 {
		t.Errorf("Level() = %v, want info", cfg.Level())
	}
	if cfg.APIURL != "https://waytutor-api.marcio.run" {
		t.Errorf("APIURL = %q, want the trailing slash removed", cfg.APIURL)
	}
	if cfg.PollEvery.String() != "5m0s" {
		t.Errorf("PollEvery = %v, want 5m", cfg.PollEvery)
	}
	if len(cfg.AdminUserIDs) != 2 {
		t.Errorf("AdminUserIDs = %v, want two ids", cfg.AdminUserIDs)
	}
}

func TestLoadRequiresEverySetting(t *testing.T) {
	keys := []string{
		"DISCORD_TOKEN",
		"DISCORD_APPLICATION_ID",
		"DISCORD_GUILD_ID",
		"DISCORD_CHANNEL_ID",
		"DISCORD_ADMIN_USER_IDS",
		"ADMIN_TOKEN",
	}

	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			requiredEnv(t)
			t.Setenv(key, "")

			if _, err := Load(); err == nil {
				t.Fatalf("Load() error = nil, want an error for a missing %s", key)
			}
		})
	}
}

func TestLoadRejectsShortPollInterval(t *testing.T) {
	requiredEnv(t)
	t.Setenv("POLL_INTERVAL", "5s")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for a poll interval under 30s")
	}
}
