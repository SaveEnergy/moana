// Package handlers implements Moana HTTP handlers: [App] holds config, store, and templates.
//
// Routing: [RegisterRoutes] delegates to routes_auth.go (login, forgot/reset password, logout), routes_dashboard.go (dashboard uses `GET /{$}` for exact `/` only), routes_ledger.go,
// routes_settings.go, routes_notifications.go (GET /notifications [App.Notifications], POST /notifications/read [App.NotificationMarkRead], notifications.html). Form field names for login/settings: formnames.go.
// Auth: [App.WithAuth] loads the user plus unread notification count in one store query; [App.CurrentUser] (e.g. login redirect) uses GetUserByID only (no unread subquery). Form parse failures: [requireParseForm] / [requireParseFormSettings] (see forms_test.go).
// Rendering: layout.go,
// transaction_*_render.go, categories_render.go. See docs/architecture.md for the full map.
//
// Router-level regression tests use package handlers_test (integration_*.go, routes_register_test.go) to avoid an import cycle;
// shared assertions live in integration_assert_test.go.
package handlers
