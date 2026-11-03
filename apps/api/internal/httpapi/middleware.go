package httpapi

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := rand.Text()
		w.Header().Set("X-Request-Id", id)

		rec := &recorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		slog.LogAttrs(r.Context(), slog.LevelInfo, "request",
			slog.String("request_id", id),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Duration("duration", time.Since(start)),
			slog.Int("bytes", rec.bytes),
			slog.String("client_ip", r.RemoteAddr),
		)
	})
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				slog.ErrorContext(r.Context(), "handler panic",
					"value", v,
					"stack", string(debug.Stack()),
				)
				writeJSON(w, http.StatusInternalServerError,
					errorEnvelope{Error: errorBody{Code: "internal_error"}},
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type recorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *recorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}
