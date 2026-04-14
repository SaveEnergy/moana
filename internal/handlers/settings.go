package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"moana/internal/auth"
	"moana/internal/household"
	"moana/internal/store"
)

func redirectSettingsErr(w http.ResponseWriter, r *http.Request, msg string) {
	http.Redirect(w, r, "/settings?err="+url.QueryEscape(msg), http.StatusSeeOther)
}

func redirectSettingsOK(w http.ResponseWriter, r *http.Request, okKey string) {
	http.Redirect(w, r, "/settings?ok="+url.QueryEscape(okKey), http.StatusSeeOther)
}

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
