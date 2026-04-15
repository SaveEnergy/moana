package handlers

import (
	"net/http"
	"strings"

	"moana/internal/auth"
	"moana/internal/store"
)

// SettingsProfileUpdate handles POST /settings/profile.
func (a *App) SettingsProfileUpdate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseFormSettings(w, r) {
		return
	}
	first := strings.TrimSpace(r.FormValue("first_name"))
	last := strings.TrimSpace(r.FormValue("last_name"))
	newPw := r.FormValue("new_password")
	ctx := r.Context()
	if err := a.Store.UpdateUserProfile(ctx, u.ID, first, last); err != nil {
		redirectSettingsErr(w, r, "Could not save profile.")
		return
	}
	if newPw != "" {
		current := r.FormValue("current_password")
		if current == "" {
			redirectSettingsErr(w, r, "Enter your current password to set a new one.")
			return
		}
		if err := auth.CheckPassword(u.PasswordHash, current); err != nil {
			redirectSettingsErr(w, r, "Current password is incorrect.")
			return
		}
		confirm := r.FormValue("new_password_confirm")
		if newPw != confirm {
			redirectSettingsErr(w, r, "New passwords do not match.")
			return
		}
		hash, err := auth.HashPassword(newPw)
		if err != nil {
			redirectSettingsErr(w, r, "Could not set password.")
			return
		}
		if err := a.Store.UpdateUserPassword(ctx, u.ID, hash); err != nil {
			redirectSettingsErr(w, r, "Could not update password.")
			return
		}
	}
	ok := "saved"
	if newPw != "" {
		ok = "password"
	}
	redirectSettingsOK(w, r, ok)
}
