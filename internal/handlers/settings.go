package handlers

import (
	"net/http"

	"moana/internal/household"
	"moana/internal/httperr"
	"moana/internal/store"
)

// Settings shows profile and household management for the signed-in user.
func (a *App) Settings(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	q := r.URL.Query()
	data, err := household.LoadSettingsPage(ctx, a.Store, u, q.Get(SettingsErrorQueryParam), q.Get(SettingsOKQueryParam))
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	a.renderShell(w, r, "settings.html", data, "Settings", "settings", "settings-shell", u)
}
