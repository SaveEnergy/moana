package app

import (
	"fmt"
	"net/http"

	"moana/internal/config"
	"moana/internal/handlers"
	"moana/internal/render"
	"moana/internal/server"
	"moana/internal/store"
	"moana/internal/tmpl"
)

// New builds an [handlers.App] with parsed HTML templates and the given config + store.
func New(cfg *config.Config, st *store.Store) (*handlers.App, error) {
	tmpl, err := tmpl.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	return &handlers.App{
		Config: cfg,
		Store:  st,
		Render: &render.Engine{Templates: tmpl},
	}, nil
}

// HTTPHandler returns the production HTTP handler (parsed templates + routes + logging).
// It passes [config.Config.MaxRequestBodyBytes] into [server.RouterOptions] so env-driven limits apply; see [server.NewRouterWithRouterOptions] for the full middleware chain.
// Tests that need a bare [handlers.App] should use [New] and [server.NewRouter] directly.
func HTTPHandler(cfg *config.Config, st *store.Store) (http.Handler, error) {
	a, err := New(cfg, st)
	if err != nil {
		return nil, err
	}
	opts := &server.RouterOptions{MaxRequestBodyBytes: cfg.MaxRequestBodyBytes}
	return server.NewRouterWithRouterOptions(opts, a), nil
}
