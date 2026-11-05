package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
	"github.com/marciojunior/waytutor/apps/bot/internal/config"
	"github.com/marciojunior/waytutor/apps/bot/internal/discord"
	"github.com/marciojunior/waytutor/apps/bot/internal/poller"
)

func main() {
	if err := run(); err != nil {
		slog.Error("bot stopped", "error", err)
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

	client := api.New(cfg.APIURL, cfg.AdminToken)

	bot, err := discord.New(cfg, discord.NewHandler(cfg, client, slog.Default()), slog.Default())
	if err != nil {
		return err
	}
	defer bot.Close()

	notifier := poller.New(client, bot.Session, cfg.ChannelID, cfg.PollEvery, slog.Default())

	// Registering on every ready makes a renamed or deleted command heal itself
	// on reconnect.
	bot.OnReady(func() {
		slog.Info("discord gateway ready")
		if err := bot.RegisterCommands(); err != nil {
			slog.Error("register commands failed", "error", err)
		}
	})

	if err := bot.Open(); err != nil {
		return err
	}

	go notifier.Run(ctx)

	slog.Info("bot started", "env", cfg.Env, "guild_id", cfg.GuildID, "poll_interval", cfg.PollEvery.String())

	<-ctx.Done()
	slog.Info("shutting down")
	return nil
}
