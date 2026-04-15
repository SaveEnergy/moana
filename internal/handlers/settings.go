package handlers

import (
	"errors"
	"net/http"

	"moana/internal/household"
	"moana/internal/store"
)

// Settings shows profile and household management for the signed-in user.
func (a *App) Settings(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	data, err := household.LoadSettingsPage(ctx, a.Store, u, r.URL.Query().Get("err"), r.URL.Query().Get("ok"))
	if err != nil {
		if errors.Is(err, household.ErrHouseholdMissing) {
			http.Error(w, "household not found", http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.renderShell(w, "settings.html", data, layoutShellMain("Settings", "settings", "settings-shell", u))
}
