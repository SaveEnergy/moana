package handlers

import (
	"errors"
	"testing"

	"moana/internal/auth"
	"moana/internal/store"
)

func TestValidateSessionUser_nilUser(t *testing.T) {
	t.Parallel()
	_, err := validateSessionUser(nil, &auth.SessionPayload{UserID: 1, Role: "user"})
	if !errors.Is(err, ErrAuthRequired) {
		t.Fatalf("got %v want %v", err, ErrAuthRequired)
	}
}

func TestValidateSessionUser_roleMismatch(t *testing.T) {
	t.Parallel()
	u := &store.User{ID: 1, Role: "user"}
	_, err := validateSessionUser(u, &auth.SessionPayload{UserID: 1, Role: "admin"})
	if !errors.Is(err, ErrAuthRequired) {
		t.Fatalf("got %v want %v", err, ErrAuthRequired)
	}
}

func TestValidateSessionUser_ok(t *testing.T) {
	t.Parallel()
	u := &store.User{ID: 1, Role: "user"}
	out, err := validateSessionUser(u, &auth.SessionPayload{UserID: 1, Role: "user"})
	if err != nil || out != u {
		t.Fatalf("got (%v, %v) want (same user, nil)", out, err)
	}
}
