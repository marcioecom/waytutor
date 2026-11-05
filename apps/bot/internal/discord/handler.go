package discord

import (
	"context"
	"errors"
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

	content := h.run(ctx, subcommand)
	if err := editResponse(session, interaction, content); err != nil {
		h.logger.Error("edit interaction response failed", "error", err)
	}
}

func (h *Handler) run(ctx context.Context, option *discordgo.ApplicationCommandInteractionDataOption) string {
	if option.Name == "pending" {
		return h.pending(ctx)
	}

	email := strings.ToLower(strings.TrimSpace(optionValue(option, "email")))
	if email == "" {
		return "Informe o e-mail da conta."
	}

	accounts, err := h.api.List(ctx, "", email)
	if err != nil {
		return h.failure("list accounts", err)
	}
	if len(accounts) == 0 {
		return fmt.Sprintf("Não encontrei conta com o e-mail %s.", email)
	}
	account := accounts[0]

	switch option.Name {
	case "approve":
		if err := h.api.SetStatus(ctx, account.ID, "active"); err != nil {
			return h.failure("approve account", err)
		}
		return fmt.Sprintf("Liberado: %s.", account.Email)

	case "pause":
		if err := h.api.SetStatus(ctx, account.ID, "paused"); err != nil {
			return h.failure("pause account", err)
		}
		return fmt.Sprintf("Pausado: %s.", account.Email)

	case "revoke":
		revoked, err := h.api.RevokeSessions(ctx, account.ID)
		if err != nil {
			return h.failure("revoke sessions", err)
		}
		return fmt.Sprintf("Sessões encerradas de %s: %d.", account.Email, revoked)
	}

	return "Comando desconhecido."
}

func (h *Handler) pending(ctx context.Context) string {
	accounts, err := h.api.List(ctx, "pending", "")
	if err != nil {
		return h.failure("list pending", err)
	}
	if len(accounts) == 0 {
		return "Nenhum cadastro aguardando."
	}

	var builder strings.Builder
	builder.WriteString("Aguardando liberação:\n")
	for _, account := range accounts {
		fmt.Fprintf(&builder, "• %s (%s)\n", account.Email, nameOr(account.Name))
	}
	builder.WriteString("\nLibere com /tutor approve email:<e-mail>")
	return builder.String()
}

func (h *Handler) failure(action string, err error) string {
	h.logger.Error("admin request failed", "action", action, "error", err)

	var status api.StatusError
	if errors.As(err, &status) && status.Code == 404 {
		return "A conta não existe mais."
	}
	return "A API recusou o pedido. Veja o log do bot."
}

func nameOr(name string) string {
	if strings.TrimSpace(name) == "" {
		return "sem nome"
	}
	return name
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
func editResponse(session *discordgo.Session, interaction *discordgo.InteractionCreate, content string) error {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral},
	})
	if err != nil {
		return fmt.Errorf("defer response: %w", err)
	}

	if _, err := session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	}); err != nil {
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
