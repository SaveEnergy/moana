package handlers

import (
	"net/http"

	"moana/internal/store"
)

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
