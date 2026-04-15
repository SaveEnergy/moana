package handlers

import "net/http"

func registerDashboardRoutes(mux *http.ServeMux, app *App) {
	// Use /{$} so only "/" matches; bare "GET /" is a prefix in Go 1.22+ ServeMux and would
	// incorrectly serve the dashboard for every path (e.g. /foo).
	mux.Handle("GET /{$}", app.WithAuth(app.Dashboard))
}
