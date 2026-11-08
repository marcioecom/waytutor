package discord

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
)

func (h *Handler) run(ctx context.Context, option *discordgo.ApplicationCommandInteractionDataOption) reply {
	switch option.Name {
	case "pending":
		return text(h.pending(ctx))
	case "usage":
		return h.usage(ctx)
	case "accounts":
		return h.accounts(ctx)
	}

	email := strings.ToLower(strings.TrimSpace(optionValue(option, "email")))
	if email == "" {
		return text("Informe o e-mail da conta.")
	}

	accounts, err := h.api.List(ctx, "", email)
	if err != nil {
		return text(h.failure("list accounts", err))
	}
	if len(accounts) == 0 {
		return text(fmt.Sprintf("Não encontrei conta com o e-mail %s.", email))
	}
	account := accounts[0]

	switch option.Name {
	case "account":
		return reply{embeds: []*discordgo.MessageEmbed{accountEmbed(account)}}

	case "approve":
		if err := h.api.SetStatus(ctx, account.ID, "active"); err != nil {
			return text(h.failure("approve account", err))
		}
		return text(fmt.Sprintf("Liberado: %s.", account.Email))

	case "pause":
		if err := h.api.SetStatus(ctx, account.ID, "paused"); err != nil {
			return text(h.failure("pause account", err))
		}
		return text(fmt.Sprintf("Pausado: %s.", account.Email))

	case "revoke":
		revoked, err := h.api.RevokeSessions(ctx, account.ID)
		if err != nil {
			return text(h.failure("revoke sessions", err))
		}
		return text(fmt.Sprintf("Sessões encerradas de %s: %d.", account.Email, revoked))
	}

	return text("Comando desconhecido.")
}

func (h *Handler) usage(ctx context.Context) reply {
	usage, err := h.api.Usage(ctx)
	if err != nil {
		return text(h.failure("load usage", err))
	}
	return reply{embeds: []*discordgo.MessageEmbed{usageEmbed(usage)}}
}

func (h *Handler) accounts(ctx context.Context) reply {
	accounts, err := h.api.List(ctx, "", "")
	if err != nil {
		return text(h.failure("list accounts", err))
	}
	return reply{embeds: []*discordgo.MessageEmbed{accountsEmbed(accounts)}}
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
