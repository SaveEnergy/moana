package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"moana/internal/auth"
	"moana/internal/household"
	"moana/internal/store"
)

// SettingsHouseholdMemberAdd handles POST /settings/household/members.
func (a *App) SettingsHouseholdMemberAdd(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseFormSettings(w, r) {
		return
	}
	if !household.CanManage(u) {
		redirectSettingsErr(w, r, "You cannot add members.")
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	if email == "" || password == "" {
		redirectSettingsErr(w, r, "Email and password are required.")
		return
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		redirectSettingsErr(w, r, "Could not set password.")
		return
	}
	ctx := r.Context()
	if _, err := a.Store.CreateHouseholdMember(ctx, u.HouseholdID, email, hash); err != nil {
		if errors.Is(err, store.ErrDuplicateUserEmail) {
			redirectSettingsErr(w, r, "A user with that email already exists.")
			return
		}
		redirectSettingsErr(w, r, "Could not add member.")
		return
	}
	redirectSettingsOK(w, r, "member")
}

// SettingsHouseholdMemberRemove handles POST /settings/household/members/remove.
func (a *App) SettingsHouseholdMemberRemove(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseFormSettings(w, r) {
		return
	}
	idStr := r.FormValue("user_id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || targetID <= 0 {
		redirectSettingsErr(w, r, "Invalid member.")
		return
	}
	ctx := r.Context()
	target, err := a.Store.GetUserByID(ctx, targetID)
	if err != nil || target == nil {
		redirectSettingsErr(w, r, "User not found.")
		return
	}
	if target.HouseholdID != u.HouseholdID {
		redirectSettingsErr(w, r, "Not in your household.")
		return
	}
	// Self-service leave
	if targetID == u.ID {
		if strings.ToLower(target.HouseholdRole) == "owner" {
			n, err := a.Store.CountHouseholdMembers(ctx, u.HouseholdID)
			if err != nil {
				redirectSettingsErr(w, r, "Could not verify household.")
				return
			}
			if n > 1 {
				redirectSettingsErr(w, r, "Transfer ownership before leaving the household.")
				return
			}
		}
		if err := a.Store.DetachUserToSoloHousehold(ctx, targetID); err != nil {
			redirectSettingsErr(w, r, "Could not leave household.")
			return
		}
		redirectSettingsOK(w, r, "left")
		return
	}
	if !household.CanRemoveMember(u, target) {
		redirectSettingsErr(w, r, "You cannot remove this member.")
		return
	}
	if err := a.Store.DetachUserToSoloHousehold(ctx, targetID); err != nil {
		redirectSettingsErr(w, r, "Could not remove member.")
		return
	}
	redirectSettingsOK(w, r, "removed")
}
