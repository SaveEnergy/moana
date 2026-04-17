package server

import (
	"log/slog"
	"net/http"
	"time"
)

// RequestLogging logs method, path, status, and duration for each request.
func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lw.status,
			"dur_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Flush implements [http.Flusher] so middleware that wraps the response (e.g.
// compression) can flush through [statusWriter] without losing the type assertion.
func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
