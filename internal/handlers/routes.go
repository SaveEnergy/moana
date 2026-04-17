package handlers

import "net/http"

// RegisterRoutes attaches all application routes to mux (static + health are registered separately by server).
func RegisterRoutes(mux *http.ServeMux, app *App) {
	registerAuthRoutes(mux, app)
	registerDashboardRoutes(mux, app)
	registerLedgerRoutes(mux, app)
	registerSettingsRoutes(mux, app)
	registerNotificationsRoutes(mux, app)
}
