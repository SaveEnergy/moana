package handlers

import "net/http"

func registerAuthRoutes(mux *http.ServeMux, app *App) {
	mux.HandleFunc(http.MethodGet+" "+LoginPath, app.LoginPage)
	mux.HandleFunc(http.MethodPost+" "+LoginPath, app.LoginSubmit)
	mux.HandleFunc(http.MethodPost+" "+LogoutPath, app.Logout)
}
