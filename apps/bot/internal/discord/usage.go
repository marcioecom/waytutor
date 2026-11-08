package discord

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
)

// The admin lives in UTC-3. Brazil dropped daylight saving in 2019, so a fixed
// offset is the whole rule and the image needs no tzdata.
var displayZone = time.FixedZone("UTC-3", -3*60*60)

const (
	colorInfo    = 0x3b82f6
	colorNeutral = 0x6b7280
)

func usageEmbed(usage api.Usage) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: "WayTutor · Uso",
		Description: fmt.Sprintf(
			"Hoje `%s` (UTC-3) · mês `%s`",
			formatDay(usage.Day), usage.Period,
		),
		Color: colorInfo,
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name: "Hoje",
			Value: fmt.Sprintf(
				"%s · %s\n%s · %s com uso",
				count(usage.Today.Requests, "áudio", "áudios"), count(usage.Today.Chars, "caractere", "caracteres"),
				count(usage.Today.Summaries, "resumo", "resumos"), count(usage.Today.Accounts, "conta", "contas"),
			),
			Inline: true,
		},
		{
			Name: "Mês",
			Value: fmt.Sprintf(
				"%s · %s\n%s da cota global · %s",
				count(usage.Month.Requests, "áudio", "áudios"), count(usage.Month.Chars, "caractere", "caracteres"),
				percent(usage.Month.Chars, usage.Limits.GlobalChars), count(usage.Month.Summaries, "resumo", "resumos"),
			),
			Inline: true,
		},
		{
			Name: "Contas",
			Value: fmt.Sprintf(
				"%s · %s\n%s",
				count(usage.Accounts.Active, "ativa", "ativas"), count(usage.Accounts.Pending, "aguardando", "aguardando"),
				count(usage.Accounts.Paused, "pausada", "pausadas"),
			),
			Inline: true,
		},
	}

	chars := make([]int64, 0, len(usage.Daily))
	for _, day := range usage.Daily {
		chars = append(chars, day.Chars)
	}
	daily := "Sem uso registrado ainda."
	if !allZero(chars) {
		daily = sparkline(chars)
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("Caracteres por dia (%s)", count(int64(len(chars)), "dia", "dias")),
		Value: daily,
	})

	if len(usage.Top) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Quem mais usou no mês",
			Value: topLines(usage.Top),
		})
	}

	return embed
}

func topLines(accounts []api.UsageAccount) string {
	lines := make([]string, 0, len(accounts))
	for _, account := range accounts {
		who := account.Email
		if name := strings.TrimSpace(account.Name); name != "" {
			who = fmt.Sprintf("%s (%s)", account.Email, truncate(name, 30))
		}
		lines = append(lines, fmt.Sprintf(
			"• %s — %s · %s",
			who, count(account.Chars, "caractere", "caracteres"), count(account.Requests, "áudio", "áudios"),
		))
	}
	return strings.Join(lines, "\n")
}

// Sparkline renders a series as block characters, the substitute for charts
// that Discord embeds cannot show.
func sparkline(values []int64) string {
	if len(values) == 0 {
		return "—"
	}

	blocks := []rune("▁▂▃▄▅▆▇█")
	maximum := int64(0)
	for _, value := range values {
		if value > maximum {
			maximum = value
		}
	}

	var builder strings.Builder
	for _, value := range values {
		index := 0
		if maximum > 0 {
			index = int(math.Round(float64(value) / float64(maximum) * float64(len(blocks)-1)))
		}
		builder.WriteRune(blocks[index])
	}
	return builder.String()
}

// count writes "1 áudio" or "2 áudios", so panel copy never reads like a bug.
func count(value int64, singular, plural string) string {
	if value == 1 {
		return number(value) + " " + singular
	}
	return number(value) + " " + plural
}

// allZero keeps an empty series from rendering as a divider made of blocks.
func allZero(values []int64) bool {
	for _, value := range values {
		if value > 0 {
			return false
		}
	}
	return true
}

func number(value int64) string {
	negative := value < 0
	digits := strconv.FormatInt(value, 10)
	if negative {
		digits = digits[1:]
	}

	var builder strings.Builder
	if negative {
		builder.WriteByte('-')
	}
	for index, digit := range digits {
		if index > 0 && (len(digits)-index)%3 == 0 {
			builder.WriteByte('.')
		}
		builder.WriteRune(digit)
	}
	return builder.String()
}

func percent(part, total int64) string {
	if total == 0 {
		return "—"
	}
	return strconv.FormatFloat(float64(part)/float64(total)*100, 'f', 0, 64) + "%"
}

func formatDay(day string) string {
	parsed, err := time.Parse("2006-01-02", day)
	if err != nil {
		return day
	}
	return parsed.Format("02/01")
}

func formatTime(value string) string {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "—"
	}
	return parsed.In(displayZone).Format("02/01 15:04")
}
