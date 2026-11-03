package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/marciojunior/waytutor/apps/api/internal/store"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code string `json:"code"`
}

func New(s *store.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		code, status := http.StatusOK, "ok"
		if err := s.Healthy(ctx); err != nil {
			code, status = http.StatusServiceUnavailable, "unavailable"
		}
		writeJSON(w, code, map[string]string{"status": status})
	})

	return recoverer(requestLog(mux))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
