package testutil

import (
	"context"
	"testing"

	"moana/internal/handlers"
	"moana/internal/passwordtest"
)

// MustCreateUser inserts a user with a bcrypt test hash (fast; see [passwordtest.MustHash]) and fails the test on error.
// Role must be "user" or "admin". Returns the new user's id.
func MustCreateUser(t *testing.T, app *handlers.App, email, password, role string) int64 {
	t.Helper()
	ctx := context.Background()
	hash := passwordtest.MustHash(t, password)
	id, err := app.Store.CreateUser(ctx, email, hash, role)
	if err != nil {
		t.Fatal(err)
	}
	return id
}
