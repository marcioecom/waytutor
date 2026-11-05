package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	BaseURL = "https://ai-gateway.vercel.sh"

	SpeechTimeout  = 120 * time.Second
	SummaryTimeout = 90 * time.Second

	maxResponseBody = 32 << 20
)

var (
	ErrAuth     = errors.New("gateway rejected the credentials")
	ErrQuota    = errors.New("gateway is out of quota")
	ErrUpstream = errors.New("gateway failed")
)

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func New(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		http:    &http.Client{},
	}
}

// The caller owns the deadline, because speech and summary time out differently.
func (c *Client) post(ctx context.Context, path string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUpstream, err)
	}
	defer res.Body.Close()

	switch {
	case res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden:
		return nil, fmt.Errorf("%w: %s", ErrAuth, res.Status)
	case res.StatusCode == http.StatusPaymentRequired || res.StatusCode == http.StatusTooManyRequests:
		return nil, fmt.Errorf("%w: %s", ErrQuota, res.Status)
	case res.StatusCode != http.StatusOK:
		detail, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return nil, fmt.Errorf("%w: %s: %s", ErrUpstream, res.Status, detail)
	}

	response, err := io.ReadAll(io.LimitReader(res.Body, maxResponseBody))
	if err != nil {
		return nil, fmt.Errorf("%w: read response: %s", ErrUpstream, err)
	}
	return response, nil
}
