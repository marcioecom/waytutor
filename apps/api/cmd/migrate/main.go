package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/marciojunior/waytutor/apps/api/internal/config"
	"github.com/marciojunior/waytutor/apps/api/internal/migrations"
)

func main() {
	if err := run(); err != nil {
		slog.Error("migrate failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	dsn, err := config.LoadDatabaseURL()
	if err != nil {
		return err
	}

	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse database url: %w", err)
	}
	// The pooler in front of Neon runs PgBouncer in transaction mode, which hands
	// a connection to another client between statements, so cached prepared
	// statements collide across them.
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db := stdlib.OpenDB(*cfg)
	defer db.Close()

	if err := apply(db); err != nil {
		return err
	}

	slog.Info("migrations applied")
	return nil
}

func apply(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
