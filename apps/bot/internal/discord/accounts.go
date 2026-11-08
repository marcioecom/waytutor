package discord

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
)

const (
	colorActive  = 0x2e7d5b
	colorPending = 0xd97706
	colorPaused  = 0x6b7280
	maxListed    = 10
)

func accountsEmbed(accounts []api.Account) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: "WayTutor · Contas",
		Color: colorNeutral,
	}

	var active, pending, paused int
	for _, account := range accounts {
		switch account.Status {
		case "active":
			active++
		case "pending":
			pending++
		case "paused":
			paused++
		}
	}
	embed.Description = fmt.Sprintf("%d ativas · %d aguardando · %d pausadas", active, pending, paused)

	if len(accounts) == 0 {
		embed.Description += "\n\nNenhuma conta cadastrada."
		return embed
	}

	sort.SliceStable(accounts, func(i, j int) bool {
		if accounts[i].UsageChars == accounts[j].UsageChars {
			return accounts[i].LastSeenAt > accounts[j].LastSeenAt
		}
		return accounts[i].UsageChars > accounts[j].UsageChars
	})

	for index, account := range accounts {
		if index == maxListed {
			embed.Description += fmt.Sprintf("\nMostrando as %d contas que mais usaram.", maxListed)
			break
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  truncate(displayName(account), 100),
			Value: fmt.Sprintf("%s · %s\n%s", account.Email, statusLabel(account.Status), usageText(account)),
		})
	}
	return embed
}

func accountEmbed(account api.Account) *discordgo.MessageEmbed {
	description := account.Email
	if name := strings.TrimSpace(account.Name); name != "" {
		description = fmt.Sprintf("**%s**\n%s", name, account.Email)
	}

	return &discordgo.MessageEmbed{
		Title:       "WayTutor · Conta",
		Description: description,
		Color:       statusColor(account.Status),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Status", Value: statusLabel(account.Status), Inline: true},
			{Name: "Criada em", Value: formatTime(account.CreatedAt), Inline: true},
			{Name: "Último acesso", Value: lastSeen(account.LastSeenAt), Inline: true},
			{Name: "Uso no mês", Value: textUsage(account), Inline: true},
			{Name: "Resumos no mês", Value: number(account.Summaries), Inline: true},
			{Name: "Média por áudio", Value: averageChars(account), Inline: true},
		},
	}
}

func usageText(account api.Account) string {
	parts := make([]string, 0, 2)
	if account.UsageRequests == 0 {
		parts = append(parts, "sem uso no mês")
	} else {
		parts = append(parts, textUsage(account))
	}
	if account.LastSeenAt == "" {
		parts = append(parts, "nunca acessou")
	} else {
		parts = append(parts, "visto "+formatTime(account.LastSeenAt))
	}
	return strings.Join(parts, " · ")
}

func textUsage(account api.Account) string {
	return fmt.Sprintf(
		"%s · %s",
		count(account.UsageChars, "caractere", "caracteres"),
		count(account.UsageRequests, "áudio", "áudios"),
	)
}

func averageChars(account api.Account) string {
	if account.UsageRequests == 0 {
		return "—"
	}
	average := account.UsageChars / account.UsageRequests
	return count(average, "caractere", "caracteres")
}

func lastSeen(value string) string {
	if value == "" {
		return "nunca"
	}
	return formatTime(value)
}

func displayName(account api.Account) string {
	if name := strings.TrimSpace(account.Name); name != "" {
		return name
	}
	return account.Email
}

func statusLabel(status string) string {
	switch status {
	case "active":
		return "ativa"
	case "pending":
		return "aguardando liberação"
	case "paused":
		return "pausada"
	default:
		return status
	}
}

func statusColor(status string) int {
	switch status {
	case "active":
		return colorActive
	case "pending":
		return colorPending
	case "paused":
		return colorPaused
	default:
		return colorNeutral
	}
}
