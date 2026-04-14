// Package handlers implements HTTP handlers for the authenticated web UI: [App] is in app.go;
// routing is in routes.go and routes_*.go; layoutShell helpers and renderShell are in layout.go;
// shared form parsing is in forms.go; HTML rendering delegates to [moana/internal/render.Engine];
// templates are parsed in [moana/internal/tmpl.Parse]. Settings are split across settings.go,
// settings_household.go (household name), and settings_members.go (invite / remove / leave).
package handlers
