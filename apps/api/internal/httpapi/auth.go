package httpapi

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/marciojunior/waytutor/apps/api/internal/auth"
	"github.com/marciojunior/waytutor/apps/api/internal/store"
)

const (
	accessSessionTTL = 7 * 24 * time.Hour
	loginCodeTTL     = 2 * time.Minute
	maxJSONBody      = 4 << 10
)

type accountBody struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type accountEnvelope struct {
	Account accountBody `json:"account"`
}

type exchangeResponse struct {
	Token   string      `json:"token"`
	Account accountBody `json:"account"`
}

func (h *api) login(w http.ResponseWriter, r *http.Request) {
	target, err := h.auth.AuthURL(r.URL.Query().Get("redirect_uri"))
	if errors.Is(err, auth.ErrInvalidRedirect) {
		writeError(w, http.StatusBadRequest, "invalid_redirect_uri")
		return
	}
	if err != nil {
		respondInternal(w, r, err)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *api) callback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	redirectURI, err := h.auth.VerifyState(query.Get("state"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_state")
		return
	}
	if query.Get("code") == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	identity, err := h.auth.Exchange(r.Context(), query.Get("code"))
	if err != nil {
		slog.ErrorContext(r.Context(), "google login failed", "error", err)
		writeError(w, http.StatusBadGateway, "oauth_failed")
		return
	}

	account, err := h.store.UpsertAccount(r.Context(), store.Identity{
		GoogleSub: identity.GoogleSub,
		Email:     identity.Email,
		Name:      identity.Name,
		Picture:   identity.Picture,
	})
	if err != nil {
		respondInternal(w, r, err)
		return
	}

	code := rand.Text()
	if err := h.store.CreateLoginCode(r.Context(), hashToken(code), account.ID, time.Now().Add(loginCodeTTL)); err != nil {
		respondInternal(w, r, err)
		return
	}

	dest, err := url.Parse(redirectURI)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_redirect_uri")
		return
	}
	params := dest.Query()
	params.Set("code", code)
	dest.RawQuery = params.Encode()
	http.Redirect(w, r, dest.String(), http.StatusFound)
}

func (h *api) exchange(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBody)

	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	account, err := h.store.ConsumeLoginCode(r.Context(), hashToken(body.Code))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusBadRequest, "invalid_code")
		return
	}
	if err != nil {
		respondInternal(w, r, err)
		return
	}

	token := rand.Text()
	expiresAt := time.Now().Add(accessSessionTTL)
	if err := h.store.CreateAccessSession(r.Context(), account.ID, hashToken(token), expiresAt, r.UserAgent()); err != nil {
		respondInternal(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, exchangeResponse{
		Token:   token,
		Account: newAccountBody(account),
	})
}

func (h *api) logout(w http.ResponseWriter, r *http.Request) {
	if token, ok := bearerToken(r); ok {
		if err := h.store.DeleteAccessSession(r.Context(), hashToken(token)); err != nil {
			respondInternal(w, r, err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *api) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, accountEnvelope{Account: newAccountBody(accountFrom(r.Context()))})
}

func (h *api) authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		account, err := h.store.AccountBySession(r.Context(), hashToken(token), time.Now().Add(accessSessionTTL))
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}
		if err != nil {
			respondInternal(w, r, err)
			return
		}

		if account.Status != store.StatusActive {
			code := "account_paused"
			if account.Status == store.StatusPending {
				code = "approval_required"
			}
			writeError(w, http.StatusForbidden, code)
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), accountKey{}, account)))
	}
}

type accountKey struct{}

func accountFrom(ctx context.Context) store.Account {
	account, _ := ctx.Value(accountKey{}).(store.Account)
	return account
}

func newAccountBody(a store.Account) accountBody {
	return accountBody{ID: a.ID, Email: a.Email, Name: a.Name, Status: a.Status}
}

func bearerToken(r *http.Request) (string, bool) {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}
	return token, true
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
