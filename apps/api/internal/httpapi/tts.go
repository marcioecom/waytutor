package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/marciojunior/waytutor/apps/api/internal/cache"
	"github.com/marciojunior/waytutor/apps/api/internal/quota"
	"github.com/marciojunior/waytutor/apps/api/internal/speech"
	"github.com/marciojunior/waytutor/apps/api/internal/store"
)

const (
	maxTTSBody  = 32 << 10
	maxTTSText  = 8000
	ttsMIMEType = "audio/wav"
)

func (h *api) tts(w http.ResponseWriter, r *http.Request) {
	account := accountFrom(r.Context())

	if !h.limiter.Allow(account.ID) {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, "rate_limited")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxTTSBody)

	var body struct {
		Text  string `json:"text"`
		Voice string `json:"voice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	text := strings.TrimSpace(body.Text)
	chars := int64(utf8.RuneCountInString(text))
	if chars == 0 {
		writeError(w, http.StatusBadRequest, "empty_text")
		return
	}
	if chars > maxTTSText {
		writeError(w, http.StatusRequestEntityTooLarge, "text_too_long")
		return
	}

	voice := strings.TrimSpace(body.Voice)
	if voice == "" {
		voice = speech.DefaultVoice
	}
	if !speech.ValidVoice(voice) {
		writeError(w, http.StatusBadRequest, "unknown_voice")
		return
	}

	key := speech.CacheKey(voice, text)

	audio, err := h.cache.Get(r.Context(), key)
	if err == nil {
		writeAudio(w, audio, key)
		return
	}
	if !errors.Is(err, cache.ErrMiss) {
		respondInternal(w, r, err)
		return
	}

	period := quota.Period(time.Now())
	if _, err := h.store.ReserveChars(r.Context(), account.ID, period, chars,
		h.limits.MonthlyCharsPerAccount, h.limits.GlobalMonthlyChars); err != nil {
		if errors.Is(err, store.ErrQuotaExceeded) {
			writeError(w, http.StatusPaymentRequired, "quota_exceeded")
			return
		}
		respondInternal(w, r, err)
		return
	}

	// A client that hangs up must not cancel audio that was already reserved.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), speech.RequestTimeout)
	defer cancel()

	audio, err = h.speech.Speak(ctx, text, voice)
	if err != nil {
		refundCtx, refundCancel := context.WithTimeout(context.WithoutCancel(r.Context()), 5*time.Second)
		defer refundCancel()
		if refundErr := h.store.RefundChars(refundCtx, account.ID, period, chars); refundErr != nil {
			slog.ErrorContext(r.Context(), "refund failed", "error", refundErr, "account_id", account.ID)
		}

		switch {
		case errors.Is(err, speech.ErrAuth):
			slog.ErrorContext(r.Context(), "speech credentials rejected", "error", err)
			writeError(w, http.StatusBadGateway, "provider_auth")
		case errors.Is(err, speech.ErrQuota):
			writeError(w, http.StatusServiceUnavailable, "provider_quota")
		default:
			slog.ErrorContext(r.Context(), "speech failed", "error", err, "account_id", account.ID)
			writeError(w, http.StatusBadGateway, "provider_error")
		}
		return
	}

	if err := h.cache.Put(r.Context(), key, audio); err != nil {
		slog.WarnContext(r.Context(), "cache write failed", "error", err, "key", key)
	}

	writeAudio(w, audio, key)
}

func writeAudio(w http.ResponseWriter, audio []byte, key string) {
	w.Header().Set("Content-Type", ttsMIMEType)
	w.Header().Set("Content-Length", strconv.Itoa(len(audio)))
	w.Header().Set("ETag", `"`+key+`"`)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	_, _ = w.Write(audio)
}
