package server

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"moana/internal/handlers"
)

// rollingMinuteLimiter records timestamps in a 60s window per key; allow returns false if the
// count would exceed max after the new hit.
type rollingMinuteLimiter struct {
	mu     sync.Mutex
	window time.Duration
	max    int
	byKey  map[string][]int64
}

func newRollingMinuteLimiter(max int, window time.Duration) *rollingMinuteLimiter {
	if max <= 0 {
		return nil
	}
	return &rollingMinuteLimiter{
		window: window,
		max:    max,
		byKey:  make(map[string][]int64),
	}
}

// allow returns true and records this hit, or false if the key is over the cap.
func (l *rollingMinuteLimiter) allow(key string) bool {
	if l == nil {
		return true
	}
	now := time.Now().UnixMilli()
	cut := now - l.window.Milliseconds()
	l.mu.Lock()
	defer l.mu.Unlock()
	ts := l.byKey[key]
	i := 0
	for i < len(ts) && ts[i] < cut {
		i++
	}
	ts = ts[i:]
	if len(ts) >= l.max {
		l.byKey[key] = ts
		return false
	}
	ts = append(ts, now)
	l.byKey[key] = ts
	return true
}

// clientIPString returns a client identifier for rate limiting. When trustForwarded is true, the
// first valid IP in X-Forwarded-For is used; otherwise RemoteAddr (strip port) is used.
func clientIPString(r *http.Request, trustForwarded bool) string {
	if trustForwarded {
		if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
			if i := strings.IndexByte(xff, ','); i >= 0 {
				xff = xff[:i]
			}
			ip := strings.TrimSpace(xff)
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// WithPostAuthRateLimit enforces per-IP rolling windows on POST /login and POST /forgot-password
// using [handlers.App.Config] limits; zero in config disables the corresponding check. This layer
// is intended to be the outermost handler before the mux (reject cheaply before large body reads).
func WithPostAuthRateLimit(app *handlers.App) func(http.Handler) http.Handler {
	var loginLim, forgotLim *rollingMinuteLimiter
	trust := false
	if app != nil && app.Config != nil {
		c := app.Config
		trust = c.TrustForwardedAddr
		loginLim = newRollingMinuteLimiter(c.RateLimitLoginPerMin, time.Minute)
		forgotLim = newRollingMinuteLimiter(c.RateLimitForgotPasswordPerMin, time.Minute)
	}
	if loginLim == nil && forgotLim == nil {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}
			ip := clientIPString(r, trust)
			switch r.URL.Path {
			case handlers.LoginPath:
				if loginLim != nil && !loginLim.allow("post_login\x00"+ip) {
					w.Header().Set("Retry-After", "60")
					http.Error(w, "Too many requests.\n", http.StatusTooManyRequests)
					return
				}
			case handlers.ForgotPasswordPath:
				if forgotLim != nil && !forgotLim.allow("post_forgot\x00"+ip) {
					w.Header().Set("Retry-After", "60")
					http.Error(w, "Too many requests.\n", http.StatusTooManyRequests)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
