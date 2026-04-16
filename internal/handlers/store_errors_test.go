package handlers

import (
	"errors"
	"fmt"
	"testing"

	"moana/internal/httperr"
	"moana/internal/store"
)

func TestUserFacingStoreMessage_invalidCategory(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrInvalidCategory)
	const want = "That category is not valid for this household."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_wrappedInvalidCategory(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("update: %w", store.ErrInvalidCategory)
	got := userFacingStoreMessage(err)
	const want = "That category is not valid for this household."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_unknownDoesNotPassthroughRawError(t *testing.T) {
	t.Parallel()
	err := errors.New("driver: SQLITE_CORRUPT (secret detail)")
	got := userFacingStoreMessage(err)
	if got != httperr.InternalMessage {
		t.Fatalf("got %q want %q", got, httperr.InternalMessage)
	}
}

func TestUserFacingStoreMessage_categoryNotFound(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrCategoryNotFound)
	const want = "That category could not be found."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_duplicateCategoryName(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrDuplicateCategoryName)
	const want = "A category with that name already exists."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_duplicateUserEmail(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrDuplicateUserEmail)
	const want = "A user with that email already exists."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_invalidUserEmail(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrInvalidUserEmail)
	const want = "That email address is not valid."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
