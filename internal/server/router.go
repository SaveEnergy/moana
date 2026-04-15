package server

import (
	"net/http"
	"time"

	"moana/internal/handlers"
)

// RouterOptions configures HTTP routing (e.g. for tests).
type RouterOptions struct {
	DisableRequestLogging bool
	// RequestTimeout, if positive, overrides [handlers.App.Config.RequestTimeout] for the request-context deadline.
	RequestTimeout time.Duration
}

func requestTimeout(opts *RouterOptions, app *handlers.App) time.Duration {
	if opts != nil && opts.RequestTimeout > 0 {
		return opts.RequestTimeout
	}
	if app != nil && app.Config != nil && app.Config.RequestTimeout > 0 {
		return app.Config.RequestTimeout
	}
	return 0
}

// NewRouter registers all production routes on mux.
func NewRouter(app *handlers.App) http.Handler {
	return NewRouterWithRouterOptions(nil, app)
}

// NewRouterWithRouterOptions registers routes with optional logging disabled (integration tests).
func NewRouterWithRouterOptions(opts *RouterOptions, app *handlers.App) http.Handler {
	mux := http.NewServeMux()
	registerStaticAndHealth(mux)
	handlers.RegisterRoutes(mux, app)

	var inner http.Handler = mux
	if d := requestTimeout(opts, app); d > 0 {
		inner = WithRequestTimeout(d)(mux)
	}
	var h http.Handler = inner
	if opts == nil || !opts.DisableRequestLogging {
		h = RequestLogging(inner)
	}
	return h
}
