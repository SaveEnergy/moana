package handlers

import "net/http"

func registerDashboardRoutes(mux *http.ServeMux, app *App) {
	mux.Handle("GET /", app.WithAuth(app.Dashboard))
}
