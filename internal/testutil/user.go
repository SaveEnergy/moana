package testutil

import (
	"context"
	"testing"

	"moana/internal/auth"
	"moana/internal/handlers"
)

// MustCreateUser inserts a user with a hashed password (Argon2id) and fails the test on error.
// Role must be "user" or "admin". Returns the new user's id.
func MustCreateUser(t *testing.T, app *handlers.App, email, password, role string) int64 {
	t.Helper()
	ctx := context.Background()
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	id, err := app.Store.CreateUser(ctx, email, hash, role)
	if err != nil {
		t.Fatal(err)
	}
	return id
}
