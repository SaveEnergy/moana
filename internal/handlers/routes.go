package handlers

import "net/http"

// RegisterRoutes attaches all application routes to mux (static + health are registered separately by server).
func RegisterRoutes(mux *http.ServeMux, app *App) {
	registerAuthRoutes(mux, app)
	registerMainRoutes(mux, app)
	registerSettingsRoutes(mux, app)
}
