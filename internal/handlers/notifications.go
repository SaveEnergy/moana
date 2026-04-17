package handlers

import (
	"net/http"

	"moana/internal/store"
)

// Notifications renders GET /notifications (empty inbox until a notification backend exists).
func (a *App) Notifications(w http.ResponseWriter, r *http.Request, u *store.User) {
	a.renderShell(w, "notifications.html", struct{}{}, layoutShell("Notifications", "notifications", u))
}
