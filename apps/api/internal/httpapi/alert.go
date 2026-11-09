package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// All identical failures inside this window collapse into one message, so an
// outage does not flood the channel.
const alertCooldown = 10 * time.Minute

// alerter posts failed requests to a Discord webhook.
type alerter struct {
	url    string
	client *http.Client
	mu     sync.Mutex
	sent   map[string]time.Time
}

func newAlerter(url string) *alerter {
	return &alerter{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
		sent:   make(map[string]time.Time),
	}
}

// serverError reports a failed request. The pattern keeps account ids out of
// the alert; an unmatched route falls back to the path.
func (a *alerter) serverError(r *http.Request, status int, requestID string, took time.Duration) {
	if a.url == "" {
		return
	}

	route := r.Pattern
	if route == "" {
		route = r.Method + " " + r.URL.Path
	}

	if !a.allow(fmt.Sprintf("%s %d", route, status)) {
		return
	}

	// The webhook must not add latency to a response that already failed.
	go a.send(fmt.Sprintf("**%d** `%s` after %s\nrequest_id `%s`",
		status, route, took.Round(time.Millisecond), requestID))
}

func (a *alerter) allow(key string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	if last, ok := a.sent[key]; ok && time.Since(last) < alertCooldown {
		return false
	}
	a.sent[key] = time.Now()
	return true
}

func (a *alerter) send(content string) {
	body, err := json.Marshal(map[string]string{"content": content})
	if err != nil {
		slog.Error("alert encode failed", "error", err)
		return
	}

	resp, err := a.client.Post(a.url, "application/json", bytes.NewReader(body))
	if err != nil {
		slog.Warn("alert failed", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		slog.Warn("alert rejected", "status", resp.StatusCode)
	}
}
