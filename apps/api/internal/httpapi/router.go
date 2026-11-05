package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/marciojunior/waytutor/apps/api/internal/auth"
	"github.com/marciojunior/waytutor/apps/api/internal/cache"
	"github.com/marciojunior/waytutor/apps/api/internal/gateway"
	"github.com/marciojunior/waytutor/apps/api/internal/ratelimit"
	"github.com/marciojunior/waytutor/apps/api/internal/store"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code string `json:"code"`
}

type Limits struct {
	MonthlyCharsPerAccount int64
	GlobalMonthlyChars     int64
	MonthlySummaries       int64
	RatePerMinute          int
}

type Deps struct {
	Store  *store.Store
	Auth   *auth.Authenticator
	Cache  cache.Cache
	Models *gateway.Client
	Limits Limits
}

type api struct {
	store   *store.Store
	auth    *auth.Authenticator
	cache   cache.Cache
	models  *gateway.Client
	limiter *ratelimit.Limiter
	limits  Limits
}

func New(d Deps) http.Handler {
	h := &api{
		store:   d.Store,
		auth:    d.Auth,
		cache:   d.Cache,
		models:  d.Models,
		limiter: ratelimit.New(d.Limits.RatePerMinute),
		limits:  d.Limits,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /auth/google", h.login)
	mux.HandleFunc("GET /auth/google/callback", h.callback)
	mux.HandleFunc("POST /v1/auth/exchange", h.exchange)
	mux.HandleFunc("POST /v1/auth/logout", h.logout)
	mux.HandleFunc("GET /v1/me", h.authenticated(h.me))
	mux.HandleFunc("POST /v1/tts", h.authenticated(h.tts))
	mux.HandleFunc("POST /v1/summary", h.authenticated(h.summary))

	return requestLog(cors(recoverer(mux)))
}

func (h *api) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	code, status := http.StatusOK, "ok"
	if err := h.store.Healthy(ctx); err != nil {
		code, status = http.StatusServiceUnavailable, "unavailable"
	}
	writeJSON(w, code, map[string]string{"status": status})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, errorEnvelope{Error: errorBody{Code: code}})
}

func respondInternal(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), "request failed", "error", err, "path", r.URL.Path)
	writeError(w, http.StatusInternalServerError, "internal_error")
}
