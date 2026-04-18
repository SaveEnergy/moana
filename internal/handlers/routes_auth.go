package handlers

import "net/http"

func registerAuthRoutes(mux *http.ServeMux, app *App) {
	mux.HandleFunc(http.MethodGet+" "+loginPath, app.LoginPage)
	mux.HandleFunc(http.MethodPost+" "+loginPath, app.LoginSubmit)
	mux.HandleFunc(http.MethodPost+" "+logoutPath, app.Logout)
}
