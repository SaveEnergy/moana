package app

import (
	"fmt"
	"net/http"

	"moana/internal/config"
	"moana/internal/handlers"
	"moana/internal/mail"
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
	var reset mail.PasswordResetSender
	if s := mail.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom); s != nil {
		reset = s
	}
	return &handlers.App{
		Config:        cfg,
		Store:         st,
		Render:        &render.Engine{Templates: tmpl},
		PasswordReset: reset,
	}, nil
}

// routerOptionsFromConfig maps env-backed settings onto [server.RouterOptions] for the production
// router. Request timeout and max body size are applied here so [server.NewRouterWithRouterOptions]
// does not rely only on [handlers.App.Config] fallbacks (both stay in sync with [config.Config]).
func routerOptionsFromConfig(cfg *config.Config) *server.RouterOptions {
	if cfg == nil {
		return nil
	}
	return &server.RouterOptions{
		MaxRequestBodyBytes: cfg.MaxRequestBodyBytes,
		RequestTimeout:      cfg.RequestTimeout,
	}
}

// HTTPHandler returns the production HTTP handler (parsed templates + routes + logging).
// It passes [config.Config] limits into [server.RouterOptions] (request deadline + POST body cap);
// a negative [config.Config.MaxRequestBodyBytes] disables the body cap (tests; [config.Load] never sets negative).
// see [server.NewRouterWithRouterOptions] for the full middleware chain.
// Tests that need a bare [handlers.App] should use [New] and [server.NewRouter] directly.
func HTTPHandler(cfg *config.Config, st *store.Store) (http.Handler, error) {
	a, err := New(cfg, st)
	if err != nil {
		return nil, err
	}
	return server.NewRouterWithRouterOptions(routerOptionsFromConfig(cfg), a), nil
}
