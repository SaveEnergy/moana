// Package handlers implements Moana HTTP handlers: [App] holds config, store, and templates.
//
// Routing: [RegisterRoutes] delegates to routes_auth.go, routes_dashboard.go (dashboard uses `GET /{$}` for exact `/` only), routes_ledger.go,
// routes_settings.go, routes_notifications.go. Auth/session: [App.CurrentUser], [App.WithAuth]. Rendering: layout.go,
// transaction_*_render.go, categories_render.go. See docs/architecture.md for the full map.
//
// Router-level regression tests use package handlers_test (integration_*.go, routes_register_test.go) to avoid an import cycle;
// shared assertions live in integration_assert_test.go.
package handlers
