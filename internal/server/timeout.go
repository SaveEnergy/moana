package server

import (
	"context"
	"net/http"
	"time"
)

// WithRequestTimeout wraps each request with a deadline via [context.WithTimeout].
// Handlers and database/sql calls using [http.Request.Context] cancel when the deadline passes.
func WithRequestTimeout(d time.Duration) func(http.Handler) http.Handler {
	if d <= 0 {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
