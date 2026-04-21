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

// readSessionPayload returns the verified session payload or [ErrAuthRequired].
func (a *App) readSessionPayload(r *http.Request) (*auth.SessionPayload, error) {
	sess, err := auth.ReadSession(r, a.Config.SessionSecret)
	if err != nil || sess == nil {
		return nil, ErrAuthRequired
	}
	return sess, nil
}

// validateSessionUser returns u when it is non-nil and matches the session role; otherwise
// [ErrAuthRequired] (deleted user or role changed server-side — treat as logout).
func validateSessionUser(u *store.User, sess *auth.SessionPayload) (*store.User, error) {
	if u == nil {
		return nil, ErrAuthRequired
	}
	if u.Role != sess.Role {
		return nil, ErrAuthRequired
	}
	return u, nil
}

// loadUserForSession resolves the session user and unread notification count in one store round trip.
// Returns [ErrAuthRequired] when there is no valid session, the user row is missing, or the role
// no longer matches the session (treat as logout). Store failures are returned unchanged.
func (a *App) loadUserForSession(r *http.Request) (*store.User, int64, error) {
	sess, err := a.readSessionPayload(r)
	if err != nil {
		return nil, 0, err
	}
	ctx := r.Context()
	u, unread, err := a.Store.GetUserByIDWithUnreadNotificationCount(ctx, sess.UserID)
	if err != nil {
		return nil, 0, err
	}
	u2, err := validateSessionUser(u, sess)
	if err != nil {
		return nil, 0, err
	}
	return u2, unread, nil
}

// CurrentUser returns the signed-in user from the session cookie, or nil with [ErrAuthRequired]
// if not authenticated. Store failures (e.g. DB down) are returned unchanged.
//
// It uses [store.Store.GetUserByID] only (no unread-count subquery). Authenticated page handlers
// use [WithAuth] + [loadUserForSession] instead, which combine user + unread count in one query.
func (a *App) CurrentUser(r *http.Request) (*store.User, error) {
	sess, err := a.readSessionPayload(r)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	u, err := a.Store.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}
	return validateSessionUser(u, sess)
}
