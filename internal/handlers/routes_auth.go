package handlers

import "net/http"

func registerAuthRoutes(mux *http.ServeMux, app *App) {
	mux.HandleFunc("GET /login", app.LoginPage)
	mux.HandleFunc("POST /login", app.LoginSubmit)
	mux.HandleFunc("POST /logout", app.Logout)
}
