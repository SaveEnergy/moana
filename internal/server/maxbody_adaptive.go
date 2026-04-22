package server

import "net/http"

// WithMaxRequestBodyBytesAdaptive applies [WithMaxRequestBodyBytes] to defaultN, except that
// requests matching uploadMethod+uploadPath use uploadMaxN (so a single route can post multipart
// bodies without raising the cap for the entire app). If defaultN is zero or negative, non-upload
// requests are not capped; uploadPath still uses uploadMaxN when it matches.
func WithMaxRequestBodyBytesAdaptive(defaultN int64, uploadMethod, uploadPath string, uploadMaxN int64) func(http.Handler) http.Handler {
	if uploadMaxN <= 0 {
		return WithMaxRequestBodyBytes(defaultN)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := defaultN
			if r.Method == uploadMethod && r.URL.Path == uploadPath {
				n = uploadMaxN
			} else if n <= 0 {
				next.ServeHTTP(w, r)
				return
			}
			if r.Body == nil {
				next.ServeHTTP(w, r)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}
