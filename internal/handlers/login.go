package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"moana/internal/auth"
	"moana/internal/httperr"
)

// LoginPage shows the sign-in form.
func (a *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	_, err := a.CurrentUser(r)
	if err == nil {
		http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
		return
	}
	if !errors.Is(err, ErrAuthRequired) {
		httperr.Internal(w, r, err)
		return
	}
	data := a.loginTemplateData("")
	if r.URL.Query().Get(LoginErrorQueryParam) != "" {
		data.Error = "Session expired or invalid. Please sign in again."
	}
	a.renderSimple(w, "login.html", data)
}

// LoginSubmit validates credentials and sets the session cookie.
func (a *App) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if !requireParseForm(w, r) {
		return
	}
	email := strings.TrimSpace(r.FormValue(LoginFieldEmail))
	password := r.FormValue(LoginFieldPassword)
	if email == "" || password == "" {
		http.Redirect(w, r, LoginPath, http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	u, err := a.Store.GetUserByEmail(ctx, email)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	if u == nil || auth.CheckPassword(u.PasswordHash, password) != nil {
		a.renderSimple(w, "login.html", a.loginTemplateData("Invalid email or password."))
		return
	}
	maxAge := a.Config.SessionMaxAge
	if r.FormValue(LoginFieldRemember) == "on" {
		maxAge = 30 * 24 * time.Hour
	}
	if err := auth.SignSession(w, a.Config.SessionSecret, auth.SessionPayload{
		UserID: u.ID,
		Role:   u.Role,
	}, maxAge, a.Config.SecureCookies); err != nil {
		httperr.Internal(w, r, err)
		return
	}
	http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
}

// Logout clears the session.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSession(w, a.Config.SecureCookies)
	http.Redirect(w, r, LoginPath, http.StatusSeeOther)
}

// loginTemplateData is passed to login.html (standalone, no app layout).
type loginTemplateData struct {
	Title   string
	Error   string
	Year    int
	RepoURL string
}

func (a *App) loginTemplateData(err string) loginTemplateData {
	return loginTemplateData{
		Title:   "Sign in",
		Error:   err,
		Year:    time.Now().UTC().Year(),
		RepoURL: a.Config.RepoURL,
	}
}
