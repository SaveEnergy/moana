package handlers

import (
	"net/http"
	"time"

	"moana/internal/auth"
)

// LoginPage shows the sign-in form.
func (a *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	if _, err := a.currentUser(r); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := struct {
		Title string
		Error string
		Year  int
	}{
		Title: "Sign in",
		Year:  time.Now().UTC().Year(),
	}
	if r.URL.Query().Get("error") != "" {
		data.Error = "Session expired or invalid. Please sign in again."
	}
	a.renderSimple(w, "login.html", data)
}

// LoginSubmit validates credentials and sets the session cookie.
func (a *App) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
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
		data := struct {
			Title string
			Error string
			Year  int
		}{
			Title: "Sign in",
			Error: "Invalid email or password.",
			Year:  time.Now().UTC().Year(),
		}
		a.renderSimple(w, "login.html", data)
		return
	}
	_ = auth.SignSession(w, a.Config.SessionSecret, auth.SessionPayload{
		UserID: u.ID,
		Role:   u.Role,
	}, a.Config.SessionMaxAge, a.Config.SecureCookies)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout clears the session.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSession(w, a.Config.SecureCookies)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *App) renderSimple(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.Templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
