package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Account struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type StatusError struct {
	Code int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("admin api answered %d", e.Code)
}

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{baseURL: baseURL, token: token, http: &http.Client{Timeout: 20 * time.Second}}
}

func (c *Client) List(ctx context.Context, status, email string) ([]Account, error) {
	query := url.Values{}
	if status != "" {
		query.Set("status", status)
	}
	if email != "" {
		query.Set("email", email)
	}

	var body struct {
		Accounts []Account `json:"accounts"`
	}
	if err := c.do(ctx, http.MethodGet, "/admin/accounts?"+query.Encode(), nil, &body); err != nil {
		return nil, err
	}
	return body.Accounts, nil
}

func (c *Client) SetStatus(ctx context.Context, id, status string) error {
	payload := map[string]string{"status": status}
	return c.do(ctx, http.MethodPost, "/admin/accounts/"+id+"/status", payload, nil)
}

func (c *Client) RevokeSessions(ctx context.Context, id string) (int64, error) {
	var body struct {
		Revoked int64 `json:"revoked"`
	}
	if err := c.do(ctx, http.MethodDelete, "/admin/accounts/"+id+"/sessions", nil, &body); err != nil {
		return 0, err
	}
	return body.Revoked, nil
}

func (c *Client) do(ctx context.Context, method, path string, payload, into any) error {
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Admin-Token", c.token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("call admin api: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return StatusError{Code: res.StatusCode}
	}
	if into == nil {
		return nil
	}

	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(into); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
