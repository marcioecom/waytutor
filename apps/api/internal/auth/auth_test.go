package auth

import (
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"
)

func testAuthenticator() *Authenticator {
	return New(Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CallbackURL:  "http://localhost:8080/auth/google/callback",
		StateSecret:  "0123456789abcdef",
	})
}

func TestAllowedRedirect(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want bool
	}{
		{"extension redirect", "https://abcdefghijklmnop.chromiumapp.org/", true},
		{"extension redirect with path", "https://abc.chromiumapp.org/callback?x=1", true},
		{"plain http", "http://abc.chromiumapp.org/", false},
		{"foreign host", "https://evil.com/", false},
		{"lookalike domain", "https://abc.chromiumapp.org.evil.com/", false},
		{"userinfo trick", "https://abc.chromiumapp.org@evil.com/", false},
		{"bare suffix", "https://chromiumapp.org/", false},
		{"not a url", "://", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowedRedirect(tt.uri); got != tt.want {
				t.Errorf("allowedRedirect(%q) = %v, want %v", tt.uri, got, tt.want)
			}
		})
	}
}

func TestAuthURLRejectsForeignRedirect(t *testing.T) {
	_, err := testAuthenticator().AuthURL("https://evil.com/")
	if !errors.Is(err, ErrInvalidRedirect) {
		t.Fatalf("AuthURL() error = %v, want ErrInvalidRedirect", err)
	}
}

func TestVerifyStateReturnsTheRequestedRedirect(t *testing.T) {
	a := testAuthenticator()
	const redirectURI = "https://abc.chromiumapp.org/"

	authURL, err := a.AuthURL(redirectURI)
	if err != nil {
		t.Fatalf("AuthURL() error = %v", err)
	}

	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("parse auth url: %v", err)
	}
	got, err := a.VerifyState(parsed.Query().Get("state"))
	if err != nil {
		t.Fatalf("VerifyState() error = %v", err)
	}
	if got != redirectURI {
		t.Errorf("VerifyState() = %q, want %q", got, redirectURI)
	}
}

func TestVerifyStateRejectsTamperedSignature(t *testing.T) {
	a := testAuthenticator()
	signed, err := a.sign(state{RedirectURI: "https://abc.chromiumapp.org/", ExpiresAt: time.Now().Add(time.Minute).Unix()})
	if err != nil {
		t.Fatalf("sign() error = %v", err)
	}

	payload, signature, _ := strings.Cut(signed, ".")
	tampered := payload + "." + signature[:len(signature)-1] + "A"

	if _, err := a.VerifyState(tampered); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("VerifyState() error = %v, want ErrInvalidState", err)
	}
}

func TestVerifyStateRejectsSignedForeignRedirect(t *testing.T) {
	a := testAuthenticator()
	signed, err := a.sign(state{RedirectURI: "https://evil.com/", ExpiresAt: time.Now().Add(time.Minute).Unix()})
	if err != nil {
		t.Fatalf("sign() error = %v", err)
	}
	if _, err := a.VerifyState(signed); !errors.Is(err, ErrInvalidRedirect) {
		t.Fatalf("VerifyState() error = %v, want ErrInvalidRedirect", err)
	}
}

func TestVerifyStateRejectsExpired(t *testing.T) {
	a := testAuthenticator()
	signed, err := a.sign(state{RedirectURI: "https://abc.chromiumapp.org/", ExpiresAt: time.Now().Add(-time.Second).Unix()})
	if err != nil {
		t.Fatalf("sign() error = %v", err)
	}
	if _, err := a.VerifyState(signed); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("VerifyState() error = %v, want ErrInvalidState", err)
	}
}
