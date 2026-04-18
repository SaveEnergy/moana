package handlers

import "net/http"

func registerNotificationsRoutes(mux *http.ServeMux, app *App) {
	mux.Handle(http.MethodGet+" "+NotificationsPath, app.WithAuth(app.Notifications))
}
