package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"moana/internal/auth"
	"moana/internal/store"
)

// AdminUsers shows user list and create form (admin only).
func (a *App) AdminUsers(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	users, err := a.Store.ListUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := adminPageData{
		Error:   r.URL.Query().Get("err"),
		Success: r.URL.Query().Get("ok"),
		Users:   users,
	}
	ld := LayoutData{
		Title:  "Users",
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: "admin",
	}
	a.renderShell(w, "admin_users_inner.html", data, ld)
}

type adminPageData struct {
	Error   string
	Success string
	Users   []store.UserSummary
}

// AdminUserCreate handles POST /admin/users (create account).
func (a *App) AdminUserCreate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Invalid form."), http.StatusSeeOther)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	role := strings.ToLower(strings.TrimSpace(r.FormValue("role")))
	tz := strings.TrimSpace(r.FormValue("timezone"))
	if email == "" || password == "" {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Email and password are required."), http.StatusSeeOther)
		return
	}
	if role != "user" && role != "admin" {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Role must be user or admin."), http.StatusSeeOther)
		return
	}
	if tz == "" {
		tz = "UTC"
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Could not set password."), http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	if _, err := a.Store.CreateUser(ctx, email, hash, role, tz); err != nil {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Could not create user (duplicate email?)."), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin/users?ok=created", http.StatusSeeOther)
}

// AdminUserPassword handles POST /admin/users/password (reset password).
func (a *App) AdminUserPassword(w http.ResponseWriter, r *http.Request, u *store.User) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Invalid form."), http.StatusSeeOther)
		return
	}
	idStr := r.FormValue("user_id")
	password := r.FormValue("password")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 || password == "" {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("User and new password are required."), http.StatusSeeOther)
		return
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Could not set password."), http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	if err := a.Store.UpdateUserPassword(ctx, id, hash); err != nil {
		http.Redirect(w, r, "/admin/users?err="+url.QueryEscape("Could not update password."), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin/users?ok=password", http.StatusSeeOther)
}