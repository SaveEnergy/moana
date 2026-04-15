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
	data, err := household.LoadSettingsPage(ctx, a.Store, u, r.URL.Query().Get("err"), r.URL.Query().Get("ok"))
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	a.renderShell(w, "settings.html", data, layoutShellMain("Settings", "settings", "settings-shell", u))
}
