package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/marciojunior/waytutor/apps/api/internal/gateway"
	"github.com/marciojunior/waytutor/apps/api/internal/quota"
	"github.com/marciojunior/waytutor/apps/api/internal/store"
)

const (
	maxSummaryBody = 512 << 10
	maxSummaryText = 120_000
)

type summaryEnvelope struct {
	Summary string `json:"summary"`
}

func (h *api) summary(w http.ResponseWriter, r *http.Request) {
	account := accountFrom(r.Context())

	if !h.limiter.Allow(account.ID) {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, "rate_limited")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxSummaryBody)

	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	text := strings.TrimSpace(body.Text)
	if text == "" {
		writeError(w, http.StatusBadRequest, "empty_text")
		return
	}
	if utf8.RuneCountInString(text) > maxSummaryText {
		writeError(w, http.StatusRequestEntityTooLarge, "text_too_long")
		return
	}

	period := quota.Period(time.Now())
	if err := h.store.ReserveSummary(r.Context(), account.ID, period, h.limits.MonthlySummaries); err != nil {
		if errors.Is(err, store.ErrQuotaExceeded) {
			writeError(w, http.StatusPaymentRequired, "quota_exceeded")
			return
		}
		respondInternal(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), gateway.SummaryTimeout)
	defer cancel()

	summary, err := h.models.Summarize(ctx, text)
	if err != nil {
		refundCtx, refundCancel := context.WithTimeout(context.WithoutCancel(r.Context()), 5*time.Second)
		defer refundCancel()
		if refundErr := h.store.RefundSummary(refundCtx, account.ID, period); refundErr != nil {
			slog.ErrorContext(r.Context(), "summary refund failed", "error", refundErr, "account_id", account.ID)
		}

		switch {
		case errors.Is(err, gateway.ErrAuth):
			slog.ErrorContext(r.Context(), "gateway credentials rejected", "error", err)
			writeError(w, http.StatusBadGateway, "provider_auth")
		case errors.Is(err, gateway.ErrQuota):
			writeError(w, http.StatusServiceUnavailable, "provider_quota")
		default:
			slog.ErrorContext(r.Context(), "summary failed", "error", err, "account_id", account.ID)
			writeError(w, http.StatusBadGateway, "provider_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, summaryEnvelope{Summary: summary})
}
