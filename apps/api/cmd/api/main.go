package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marciojunior/waytutor/apps/api/internal/auth"
	"github.com/marciojunior/waytutor/apps/api/internal/cache"
	"github.com/marciojunior/waytutor/apps/api/internal/config"
	"github.com/marciojunior/waytutor/apps/api/internal/gateway"
	"github.com/marciojunior/waytutor/apps/api/internal/httpapi"
	"github.com/marciojunior/waytutor/apps/api/internal/store"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Level()})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := store.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	authenticator := auth.New(auth.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		CallbackURL:  cfg.BaseURL + "/auth/google/callback",
		StateSecret:  cfg.OAuthStateSecret,
	})

	audioCache, err := newCache(cfg)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr: cfg.Addr,
		Handler: httpapi.New(httpapi.Deps{
			Store:  db,
			Auth:   authenticator,
			Cache:  audioCache,
			Models: gateway.New(cfg.GatewayAPIKey, gateway.BaseURL),
			Limits: httpapi.Limits{
				MonthlyCharsPerAccount: cfg.MonthlyCharsPerAccount,
				GlobalMonthlyChars:     cfg.GlobalMonthlyChars,
				MonthlySummaries:       cfg.MonthlySummaries,
				RatePerMinute:          cfg.RatePerMinute,
			},
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		slog.Info("api listening", "addr", cfg.Addr, "env", cfg.Env)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	return nil
}

// Production has no volume, so it caches in R2; a laptop caches on disk.
func newCache(cfg config.Config) (cache.Cache, error) {
	if cfg.R2Enabled() {
		return cache.NewR2(cache.R2Config{
			Endpoint:        cfg.R2Endpoint,
			Bucket:          cfg.R2Bucket,
			AccessKeyID:     cfg.R2AccessKeyID,
			SecretAccessKey: cfg.R2SecretAccessKey,
			Region:          cfg.R2Region,
		}), nil
	}

	disk, err := cache.NewDisk(cfg.CacheDir)
	if err != nil {
		return nil, err
	}
	slog.Warn("R2 is not configured, caching audio on disk", "dir", cfg.CacheDir)
	return disk, nil
}
