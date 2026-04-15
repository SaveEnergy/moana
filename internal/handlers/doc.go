// Package handlers implements HTTP handlers for the authenticated web UI: [App] is in app.go;
// routing is in routes.go and routes_*.go (auth, dashboard, ledger, settings); CurrentUser in current_user.go,
// WithAuth in middleware_auth.go; sign-in flow in login.go; layoutShell helpers and renderShell are in layout.go;
// shared form parsing is in forms.go; path segment helpers in params.go; transaction flows split
// across transaction_create.go, transaction_form_types.go, transaction_new_render.go, transaction_edit.go, transaction_edit_render.go;
// category error re-render in categories_render.go; HTML rendering delegates to [moana/internal/render.Engine];
// templates are parsed in [moana/internal/tmpl.Parse]. Settings are split across settings.go,
// settings_household.go (household name), settings_members.go (invite / remove / leave),
// settings_profile.go (name + password), and settings_redirect.go (flash redirects).
package handlers
