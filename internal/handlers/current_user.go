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

// CurrentUser returns the signed-in user from the session cookie, or nil with [ErrAuthRequired]
// if not authenticated. Store failures (e.g. DB down) are returned unchanged.
func (a *App) CurrentUser(r *http.Request) (*store.User, error) {
	sess, err := auth.ReadSession(r, a.Config.SessionSecret)
	if err != nil || sess == nil {
		return nil, ErrAuthRequired
	}
	ctx := r.Context()
	u, err := a.Store.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrAuthRequired
	}
	if u.Role != sess.Role {
		// role changed server-side; treat as logout
		return nil, ErrAuthRequired
	}
	return u, nil
}
