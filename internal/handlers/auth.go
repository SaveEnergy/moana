package handlers

import (
	"fmt"
	"net/http"
	"time"

	"moana/internal/auth"
	"moana/internal/store"
)

// CurrentUser returns the signed-in user from the session cookie, or nil / error if not authenticated.
func (a *App) CurrentUser(r *http.Request) (*store.User, error) {
	sess, err := auth.ReadSession(r, a.Config.SessionSecret)
	if err != nil || sess == nil {
		return nil, err
	}
	ctx := r.Context()
	u, err := a.Store.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("user missing")
	}
	if u.Role != sess.Role {
		// role changed server-side; treat as logout
		return nil, fmt.Errorf("stale session")
	}
	return u, nil
}

// WithAuth requires a valid session and loads the current user.
func (a *App) WithAuth(next func(http.ResponseWriter, *http.Request, *store.User)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := a.CurrentUser(r)
		if err != nil || u == nil {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}
		next(w, r, u)
	})
}

// LoginPage shows the sign-in form.
func (a *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	if _, err := a.CurrentUser(r); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := a.loginTemplateData("")
	if r.URL.Query().Get("error") != "" {
		data.Error = "Session expired or invalid. Please sign in again."
	}
	a.renderSimple(w, "login.html", data)
}

// LoginSubmit validates credentials and sets the session cookie.
func (a *App) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if !requireParseForm(w, r) {
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	u, err := a.Store.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if u == nil || auth.CheckPassword(u.PasswordHash, password) != nil {
		a.renderSimple(w, "login.html", a.loginTemplateData("Invalid email or password."))
		return
	}
	maxAge := a.Config.SessionMaxAge
	if r.FormValue("remember") == "on" {
		maxAge = 30 * 24 * time.Hour
	}
	_ = auth.SignSession(w, a.Config.SessionSecret, auth.SessionPayload{
		UserID: u.ID,
		Role:   u.Role,
	}, maxAge, a.Config.SecureCookies)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout clears the session.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSession(w, a.Config.SecureCookies)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
