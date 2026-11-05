package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
)

const maxChoices = 25

// approve offers what is waiting, pause offers what is in use, and revoke offers
// any status because a paused account still has sessions to end.
func autocompleteStatus(subcommand string) string {
	switch subcommand {
	case "approve":
		return "pending"
	case "pause":
		return "active"
	default:
		return ""
	}
}

func emailChoices(accounts []api.Account, query string) []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, maxChoices)
	needle := strings.ToLower(strings.TrimSpace(query))

	for _, account := range accounts {
		if needle != "" && !strings.Contains(strings.ToLower(account.Email), needle) {
			continue
		}

		label := account.Email
		if name := strings.TrimSpace(account.Name); name != "" {
			label = account.Email + " (" + name + ")"
		}

		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  truncate(label, 100),
			Value: truncate(account.Email, 100),
		})
		if len(choices) == maxChoices {
			break
		}
	}
	return choices
}

func truncate(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}
