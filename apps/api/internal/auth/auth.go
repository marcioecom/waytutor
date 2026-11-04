package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var (
	ErrInvalidState    = errors.New("invalid oauth state")
	ErrInvalidRedirect = errors.New("redirect uri is not allowed")
	ErrUnverifiedEmail = errors.New("google account email is not verified")
)

const (
	stateTTL        = 10 * time.Minute
	userInfoURL     = "https://openidconnect.googleapis.com/v1/userinfo"
	maxUserInfoBody = 64 << 10
)

type Config struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	StateSecret  string
}

type Identity struct {
	GoogleSub string
	Email     string
	Name      string
	Picture   string
}

type Authenticator struct {
	oauth  *oauth2.Config
	secret []byte
	client *http.Client
}

func New(cfg Config) *Authenticator {
	return &Authenticator{
		oauth: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.CallbackURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
			Scopes: []string{"openid", "email", "profile"},
		},
		secret: []byte(cfg.StateSecret),
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (a *Authenticator) AuthURL(redirectURI string) (string, error) {
	if !allowedRedirect(redirectURI) {
		return "", ErrInvalidRedirect
	}
	signed, err := a.sign(state{
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().Add(stateTTL).Unix(),
	})
	if err != nil {
		return "", err
	}
	return a.oauth.AuthCodeURL(signed), nil
}

func (a *Authenticator) VerifyState(signed string) (string, error) {
	st, err := a.verify(signed)
	if err != nil {
		return "", err
	}
	if !allowedRedirect(st.RedirectURI) {
		return "", ErrInvalidRedirect
	}
	return st.RedirectURI, nil
}

func (a *Authenticator) Exchange(ctx context.Context, code string) (Identity, error) {
	token, err := a.oauth.Exchange(ctx, code)
	if err != nil {
		return Identity{}, fmt.Errorf("exchange google code: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return Identity{}, fmt.Errorf("build userinfo request: %w", err)
	}
	token.SetAuthHeader(req)

	resp, err := a.client.Do(req)
	if err != nil {
		return Identity{}, fmt.Errorf("fetch google userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return Identity{}, fmt.Errorf("google userinfo: %s: %s", resp.Status, body)
	}

	var info struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxUserInfoBody)).Decode(&info); err != nil {
		return Identity{}, fmt.Errorf("decode google userinfo: %w", err)
	}
	if info.Sub == "" || info.Email == "" {
		return Identity{}, errors.New("google userinfo: missing sub or email")
	}
	if !info.EmailVerified {
		return Identity{}, ErrUnverifiedEmail
	}

	return Identity{
		GoogleSub: info.Sub,
		Email:     info.Email,
		Name:      info.Name,
		Picture:   info.Picture,
	}, nil
}

type state struct {
	RedirectURI string `json:"r"`
	ExpiresAt   int64  `json:"e"`
}

func (a *Authenticator) sign(st state) (string, error) {
	payload, err := json.Marshal(st)
	if err != nil {
		return "", fmt.Errorf("encode oauth state: %w", err)
	}
	return encode(payload) + "." + encode(a.mac(payload)), nil
}

func (a *Authenticator) verify(signed string) (state, error) {
	payloadPart, signaturePart, ok := strings.Cut(signed, ".")
	if !ok {
		return state{}, ErrInvalidState
	}
	payload, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return state{}, ErrInvalidState
	}
	signature, err := base64.RawURLEncoding.DecodeString(signaturePart)
	if err != nil {
		return state{}, ErrInvalidState
	}
	if !hmac.Equal(signature, a.mac(payload)) {
		return state{}, ErrInvalidState
	}

	var st state
	if err := json.Unmarshal(payload, &st); err != nil {
		return state{}, ErrInvalidState
	}
	if time.Now().Unix() > st.ExpiresAt {
		return state{}, ErrInvalidState
	}
	return st, nil
}

func (a *Authenticator) mac(payload []byte) []byte {
	mac := hmac.New(sha256.New, a.secret)
	mac.Write(payload)
	return mac.Sum(nil)
}

func encode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func allowedRedirect(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil {
		return false
	}
	const suffix = ".chromiumapp.org"
	host := strings.ToLower(u.Hostname())
	return len(host) > len(suffix) && strings.HasSuffix(host, suffix)
}
