package httpapi

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// cors is a wildcard because the extension origin changes with every unpacked
// install, and the session travels in a header instead of a cookie.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// The health check, the panel and the bot poll these, so logging each healthy
// call repeats one line every few seconds. A failure still logs.
var quietRoutes = map[string]bool{
	"GET /healthz":        true,
	"GET /v1/me":          true,
	"GET /admin/accounts": true,
}

func requestLog(next http.Handler, alerts *alerter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := rand.Text()
		w.Header().Set("X-Request-Id", id)

		rec := &recorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		took := time.Since(start)
		if rec.status >= 500 {
			alerts.serverError(r, rec.status, id, took)
		}
		if rec.status < 400 && quietRoutes[r.Pattern] {
			return
		}

		slog.LogAttrs(r.Context(), slog.LevelInfo, "request",
			slog.String("request_id", id),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Duration("duration", took),
			slog.Int("bytes", rec.bytes),
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
