package discord

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
	"github.com/marciojunior/waytutor/apps/bot/internal/config"
)

const requestTimeout = 20 * time.Second

type Handler struct {
	api    *api.Client
	logger *slog.Logger
	admins map[string]bool
}

func NewHandler(cfg config.Config, client *api.Client, logger *slog.Logger) *Handler {
	admins := make(map[string]bool, len(cfg.AdminUserIDs))
	for _, id := range cfg.AdminUserIDs {
		admins[strings.TrimSpace(id)] = true
	}

	return &Handler{api: client, logger: logger, admins: admins}
}

type reply struct {
	content string
	embeds  []*discordgo.MessageEmbed
}

func text(content string) reply {
	return reply{content: content}
}

func (h *Handler) OnInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		h.handleCommand(session, interaction)
	case discordgo.InteractionApplicationCommandAutocomplete:
		h.handleAutocomplete(session, interaction)
	}
}

func (h *Handler) handleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ApplicationCommandData()
	if data.Name != "tutor" || len(data.Options) == 0 {
		return
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			h.logger.Error("panic while handling interaction", "panic", recovered)
		}
	}()

	if !h.admins[userID(interaction)] {
		h.reply(session, interaction, "Sem permissão.")
		return
	}

	subcommand := data.Options[0]
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	answer := h.run(ctx, subcommand)
	if err := editResponse(session, interaction, answer); err != nil {
		h.logger.Error("edit interaction response failed", "error", err)
	}
}

func userID(interaction *discordgo.InteractionCreate) string {
	if interaction.Member != nil && interaction.Member.User != nil {
		return interaction.Member.User.ID
	}
	if interaction.User != nil {
		return interaction.User.ID
	}
	return ""
}

func optionValue(option *discordgo.ApplicationCommandInteractionDataOption, name string) string {
	for _, candidate := range option.Options {
		if candidate.Name == name {
			value, _ := candidate.Value.(string)
			return value
		}
	}
	return ""
}

func (h *Handler) reply(session *discordgo.Session, interaction *discordgo.InteractionCreate, content string) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		h.logger.Error("respond failed", "error", err)
	}
}

// The API can be cold, so answer within Discord's three seconds and edit later.
func editResponse(session *discordgo.Session, interaction *discordgo.InteractionCreate, answer reply) error {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral},
	})
	if err != nil {
		return fmt.Errorf("defer response: %w", err)
	}

	edit := &discordgo.WebhookEdit{}
	if answer.content != "" {
		edit.Content = &answer.content
	}
	if len(answer.embeds) > 0 {
		edit.Embeds = &answer.embeds
	}

	if _, err := session.InteractionResponseEdit(interaction.Interaction, edit); err != nil {
		return fmt.Errorf("edit response: %w", err)
	}
	return nil
}

func (h *Handler) handleAutocomplete(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if !h.admins[userID(interaction)] {
		h.respondAutocomplete(session, interaction, nil)
		return
	}

	data := interaction.ApplicationCommandData()
	if data.Name != "tutor" || len(data.Options) == 0 {
		h.respondAutocomplete(session, interaction, nil)
		return
	}

	subcommand := data.Options[0]

	// Discord drops an autocomplete that takes longer than three seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	accounts, err := h.api.List(ctx, autocompleteStatus(subcommand.Name), "")
	if err != nil {
		h.logger.Warn("autocomplete failed", "subcommand", subcommand.Name, "error", err)
		h.respondAutocomplete(session, interaction, nil)
		return
	}

	h.respondAutocomplete(session, interaction, emailChoices(accounts, optionValue(subcommand, "email")))
}

func (h *Handler) respondAutocomplete(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	choices []*discordgo.ApplicationCommandOptionChoice,
) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{Choices: choices},
	})
	if err != nil {
		h.logger.Warn("respond autocomplete failed", "error", err)
	}
}
