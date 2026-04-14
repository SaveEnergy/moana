package handlers

import (
	"net/http"
	"strings"

	"moana/internal/household"
	"moana/internal/store"
)

// SettingsHouseholdUpdate handles POST /settings/household (name).
func (a *App) SettingsHouseholdUpdate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseFormSettings(w, r) {
		return
	}
	name := strings.TrimSpace(r.FormValue("household_name"))
	if name == "" {
		redirectSettingsErr(w, r, "Household name is required.")
		return
	}
	if !household.CanManage(u) {
		redirectSettingsErr(w, r, "You cannot edit this household.")
		return
	}
	ctx := r.Context()
	if err := a.Store.UpdateHouseholdName(ctx, u.HouseholdID, name); err != nil {
		redirectSettingsErr(w, r, "Could not save household name.")
		return
	}
	redirectSettingsOK(w, r, "household")
}
