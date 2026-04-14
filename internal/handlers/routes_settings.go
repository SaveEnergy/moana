package handlers

import "net/http"

func registerSettingsRoutes(mux *http.ServeMux, app *App) {
	mux.Handle("GET /settings", app.WithAuth(app.Settings))
	mux.Handle("POST /settings/profile", app.WithAuth(app.SettingsProfileUpdate))
	mux.Handle("POST /settings/household", app.WithAuth(app.SettingsHouseholdUpdate))
	mux.Handle("POST /settings/household/members", app.WithAuth(app.SettingsHouseholdMemberAdd))
	mux.Handle("POST /settings/household/members/remove", app.WithAuth(app.SettingsHouseholdMemberRemove))
}
