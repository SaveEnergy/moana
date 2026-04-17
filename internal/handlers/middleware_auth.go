package handlers

import (
	"errors"
	"net/http"

	"moana/internal/httperr"
	"moana/internal/store"
)

// WithAuth requires a valid session and loads the current user.
func (a *App) WithAuth(next func(http.ResponseWriter, *http.Request, *store.User)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := a.CurrentUser(r)
		if err != nil {
			if errors.Is(err, ErrAuthRequired) {
				http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
				return
			}
			httperr.Internal(w, r, err)
			return
		}
		next(w, r, u)
	})
}
