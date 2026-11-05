package poller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
)

type Poller struct {
	client  *api.Client
	session *discordgo.Session
	channel string
	every   time.Duration
	logger  *slog.Logger
	seen    map[string]bool
}

func New(client *api.Client, session *discordgo.Session, channel string, every time.Duration, logger *slog.Logger) *Poller {
	return &Poller{
		client:  client,
		session: session,
		channel: channel,
		every:   every,
		logger:  logger,
		seen:    make(map[string]bool),
	}
}

func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.every)
	defer ticker.Stop()

	p.check(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.check(ctx)
		}
	}
}

func (p *Poller) check(ctx context.Context) {
	pending, err := p.client.List(ctx, "pending", "")
	if err != nil {
		p.logger.Error("poll failed", "error", err)
		return
	}

	fresh := make([]api.Account, 0, len(pending))
	for _, account := range pending {
		if !p.seen[account.ID] {
			fresh = append(fresh, account)
		}
	}
	if len(fresh) == 0 {
		return
	}
	for _, account := range fresh {
		p.seen[account.ID] = true
	}

	if _, err := p.session.ChannelMessageSend(p.channel, format(fresh)); err != nil {
		p.logger.Error("notify failed", "error", err)
	}
}

// A restart forgets what was announced, so the first poll sends every pending
// account in one message instead of one message per account.
func format(accounts []api.Account) string {
	var builder strings.Builder
	if len(accounts) == 1 {
		builder.WriteString("Cadastro aguardando liberação:\n")
	} else {
		fmt.Fprintf(&builder, "%d cadastros aguardando liberação:\n", len(accounts))
	}

	for _, account := range accounts {
		name := strings.TrimSpace(account.Name)
		if name == "" {
			name = "sem nome"
		}
		fmt.Fprintf(&builder, "• %s (%s)\n", account.Email, name)
	}

	builder.WriteString("\nLibere com /tutor approve email:<e-mail>")
	return builder.String()
}
