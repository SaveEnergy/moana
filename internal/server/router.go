package server

import (
	"net/http"

	"moana/internal/handlers"
)

// RouterOptions configures HTTP routing (e.g. for tests).
type RouterOptions struct {
	DisableRequestLogging bool
}

// NewRouter registers all application routes on mux.
func NewRouter(app *handlers.App) http.Handler {
	return NewRouterWithRouterOptions(nil, app)
}

// NewRouterWithRouterOptions registers routes with optional logging disabled (integration tests).
func NewRouterWithRouterOptions(opts *RouterOptions, app *handlers.App) http.Handler {
	mux := http.NewServeMux()
	registerStaticAndHealth(mux)
	handlers.RegisterRoutes(mux, app)

	var h http.Handler = mux
	if opts == nil || !opts.DisableRequestLogging {
		h = RequestLogging(mux)
	}
	return h
}
