package httpapi

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/marciojunior/waytutor/apps/api/internal/quota"
	"github.com/marciojunior/waytutor/apps/api/internal/store"
)

const (
	maxAdminBody   = 1 << 10
	adminListLimit = 200
)

var allowedStatuses = map[string]bool{
	store.StatusPending: true,
	store.StatusActive:  true,
	store.StatusPaused:  true,
}

type adminAccountBody struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
	UsageChars    int64  `json:"usageChars"`
	UsageRequests int64  `json:"usageRequests"`
	Summaries     int64  `json:"summaries"`
	LastSeenAt    string `json:"lastSeenAt,omitempty"`
}

func (h *api) admin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("X-Admin-Token")
		if subtle.ConstantTimeCompare([]byte(got), []byte(h.adminToken)) != 1 {
			writeError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}
		next(w, r)
	}
}

func (h *api) listAccounts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	period := quota.Period(time.Now())
	accounts, err := h.store.ListAccounts(r.Context(), query.Get("status"), query.Get("email"), period, adminListLimit)
	if err != nil {
		respondInternal(w, r, err)
		return
	}

	body := make([]adminAccountBody, 0, len(accounts))
	for _, account := range accounts {
		item := adminAccountBody{
			ID:            account.ID,
			Email:         account.Email,
			Name:          account.Name,
			Status:        account.Status,
			CreatedAt:     account.CreatedAt.UTC().Format(time.RFC3339),
			UsageChars:    account.UsageChars,
			UsageRequests: account.UsageRequests,
			Summaries:     account.Summaries,
		}
		if account.LastSeenAt != nil {
			item.LastSeenAt = account.LastSeenAt.UTC().Format(time.RFC3339)
		}
		body = append(body, item)
	}
	writeJSON(w, http.StatusOK, map[string]any{"accounts": body})
}

func (h *api) setAccountStatus(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBody)

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	status := strings.TrimSpace(body.Status)
	if !allowedStatuses[status] {
		writeError(w, http.StatusBadRequest, "invalid_status")
		return
	}

	err := h.store.SetAccountStatus(r.Context(), r.PathValue("id"), status)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	if err != nil {
		respondInternal(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *api) revokeSessions(w http.ResponseWriter, r *http.Request) {
	revoked, err := h.store.RevokeSessions(r.Context(), r.PathValue("id"))
	if err != nil {
		respondInternal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"revoked": revoked})
}
