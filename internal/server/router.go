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
	// MaxRequestBodyBytes, if positive, overrides [defaultMaxRequestBodyBytes] for [WithMaxRequestBodyBytes].
	// If negative, request body size is not capped (for tests; not recommended for production).
	// Zero uses the default 1 MiB cap.
	MaxRequestBodyBytes int64
}

// NewRouter registers all production routes on mux.
func NewRouter(app *handlers.App) http.Handler {
	return NewRouterWithRouterOptions(nil, app)
}

// NewRouterWithRouterOptions registers routes with optional logging disabled (integration tests).
//
// Middleware order (outermost first): [RequestLogging] unless [RouterOptions.DisableRequestLogging];
// then [wrapRouterMiddleware] — [WithRequestTimeout] when a positive deadline is configured
// ([RouterOptions.RequestTimeout] or [handlers.App.Config.RequestTimeout]), wrapping
// [WithMaxRequestBodyBytes] from [maxRequestBodyBytes] (default 1 MiB unless overridden), wrapping
// the [http.ServeMux] (static, health, [handlers.RegisterRoutes]).
func NewRouterWithRouterOptions(opts *RouterOptions, app *handlers.App) http.Handler {
	mux := http.NewServeMux()
	registerStaticAndHealth(mux)
	handlers.RegisterRoutes(mux, app)

	h := wrapRouterMiddleware(mux, opts, app)
	if opts == nil || !opts.DisableRequestLogging {
		h = RequestLogging(h)
	}
	return h
}
