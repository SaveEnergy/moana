package handlers

import (
	"errors"
	"net/http"

	"moana/internal/httperr"
	"moana/internal/store"
)

// WithAuth requires a valid session and loads the current user. It attaches the unread notification
// count to the request context for [App.renderShell] so the shell does not run a second COUNT query.
func (a *App) WithAuth(next func(http.ResponseWriter, *http.Request, *store.User)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, unread, err := a.loadUserForSession(r)
		if err != nil {
			if errors.Is(err, ErrAuthRequired) {
				http.Redirect(w, r, LoginRedirectAuth, http.StatusSeeOther)
				return
			}
			httperr.Internal(w, r, err)
			return
		}
		next(w, withUnreadNotificationCount(r, unread), u)
	})
}
