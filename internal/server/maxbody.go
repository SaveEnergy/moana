package server

import "net/http"

// defaultMaxRequestBodyBytes caps URL-encoded form POST bodies. The net/http parser allows up to
// 10 MiB per request; this keeps memory bounded for typical app forms.
const defaultMaxRequestBodyBytes = 1 << 20 // 1 MiB

// WithMaxRequestBodyBytes wraps r.Body with [http.MaxBytesReader]. If n <= 0, next is returned unchanged.
func WithMaxRequestBodyBytes(n int64) func(http.Handler) http.Handler {
	if n <= 0 {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, n)
			}
			next.ServeHTTP(w, r)
		})
	}
}
