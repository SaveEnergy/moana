package handlers

import "net/http"

func registerSettingsRoutes(mux *http.ServeMux, app *App) {
	mux.Handle(http.MethodGet+" "+SettingsPath, app.WithAuth(app.Settings))
	mux.Handle(http.MethodPost+" "+SettingsProfilePath, app.WithAuth(app.SettingsProfileUpdate))
	mux.Handle(http.MethodPost+" "+SettingsHouseholdPath, app.WithAuth(app.SettingsHouseholdUpdate))
	mux.Handle(http.MethodPost+" "+SettingsHouseholdMembersPath, app.WithAuth(app.SettingsHouseholdMemberAdd))
	mux.Handle(http.MethodPost+" "+SettingsHouseholdMembersRemovePath, app.WithAuth(app.SettingsHouseholdMemberRemove))
}
