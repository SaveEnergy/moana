package handlers

import "net/http"

func registerNotificationsRoutes(mux *http.ServeMux, app *App) {
	mux.Handle("GET /notifications", app.WithAuth(app.Notifications))
}
