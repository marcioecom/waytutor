package discord

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/marciojunior/waytutor/apps/bot/internal/config"
)

type Bot struct {
	Session *discordgo.Session
	cfg     config.Config
	logger  *slog.Logger
}

func New(cfg config.Config, handler *Handler, logger *slog.Logger) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}

	// Slash commands only need the guilds intent. The bot never reads message
	// content, so no privileged intent has to be enabled in the portal.
	session.Identify.Intents = discordgo.IntentsGuilds
	session.ShouldReconnectOnError = true
	session.AddHandler(handler.OnInteraction)

	return &Bot{Session: session, cfg: cfg, logger: logger}, nil
}

func (b *Bot) Open() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("open discord gateway: %w", err)
	}
	return nil
}

func (b *Bot) Close() error {
	return b.Session.Close()
}

// OnReady runs after every successful gateway connection, including reconnects.
func (b *Bot) OnReady(callback func()) {
	b.Session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		callback()
	})
}

func (b *Bot) RegisterCommands() error {
	_, err := b.Session.ApplicationCommandBulkOverwrite(
		b.cfg.ApplicationID,
		b.cfg.GuildID,
		[]*discordgo.ApplicationCommand{tutorCommand()},
	)
	if err != nil {
		return fmt.Errorf("register guild commands: %w", err)
	}
	b.logger.Info("discord commands registered", "guild_id", b.cfg.GuildID)
	return nil
}
