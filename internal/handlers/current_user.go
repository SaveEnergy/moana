package handlers

import (
	"errors"
	"net/http"

	"moana/internal/auth"
	"moana/internal/store"
)

// ErrAuthRequired is returned by [App.CurrentUser] when there is no valid session or the
// account cannot be resolved (deleted user, role mismatch). Callers should redirect to
// login. Database errors from the store are returned as-is so callers can respond with 500.
var ErrAuthRequired = errors.New("authentication required")

// loadUserForSession resolves the session user and unread notification count in one store round trip.
// Returns [ErrAuthRequired] when there is no valid session, the user row is missing, or the role
// no longer matches the session (treat as logout). Store failures are returned unchanged.
func (a *App) loadUserForSession(r *http.Request) (*store.User, int64, error) {
	sess, err := auth.ReadSession(r, a.Config.SessionSecret)
	if err != nil || sess == nil {
		return nil, 0, ErrAuthRequired
	}
	ctx := r.Context()
	u, unread, err := a.Store.GetUserByIDWithUnreadNotificationCount(ctx, sess.UserID)
	if err != nil {
		return nil, 0, err
	}
	if u == nil {
		return nil, 0, ErrAuthRequired
	}
	if u.Role != sess.Role {
		// role changed server-side; treat as logout
		return nil, 0, ErrAuthRequired
	}
	return u, unread, nil
}

// CurrentUser returns the signed-in user from the session cookie, or nil with [ErrAuthRequired]
// if not authenticated. Store failures (e.g. DB down) are returned unchanged.
func (a *App) CurrentUser(r *http.Request) (*store.User, error) {
	u, _, err := a.loadUserForSession(r)
	return u, err
}
