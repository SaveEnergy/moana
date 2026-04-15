package handlers

import (
	"fmt"
	"net/http"

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
