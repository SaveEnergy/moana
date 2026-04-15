// Package handlers implements Moana HTTP handlers: [App] holds config, store, and templates.
//
// Routing: [RegisterRoutes] delegates to routes_auth.go, routes_dashboard.go, routes_ledger.go,
// routes_settings.go. Auth/session: [App.CurrentUser], [App.WithAuth]. Rendering: layout.go,
// transaction_*_render.go, categories_render.go. See docs/architecture.md for the full map.
//
// Router-level regression tests use package handlers_test (integration_*.go) to avoid an import cycle.
package handlers
