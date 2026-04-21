package handlers

import (
	"database/sql"
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

func TestUserFacingStoreMessage_userNotInHousehold(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrUserNotInHousehold)
	const want = "That account cannot record transactions in this household."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_wrappedUserNotInHousehold(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("create: %w", store.ErrUserNotInHousehold)
	got := userFacingStoreMessage(err)
	const want = "That account cannot record transactions in this household."
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

func TestUserFacingStoreMessage_sqlErrNoRows(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(sql.ErrNoRows)
	const want = "That record could not be found."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_wrappedSqlErrNoRows(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("update: %w", sql.ErrNoRows)
	got := userFacingStoreMessage(err)
	const want = "That record could not be found."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
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

func TestUserFacingStoreMessage_notificationNotFound(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrNotificationNotFound)
	const want = "That notification could not be found."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_invalidNotificationBody(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(store.ErrInvalidNotificationBody)
	const want = "Notification text cannot be empty."
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestUserFacingStoreMessage_invalidCategoryAmountKind_mapsToInternal(t *testing.T) {
	t.Parallel()
	// Used by dashboard aggregates only; handlers must not surface raw strings if this ever leaks.
	got := userFacingStoreMessage(store.ErrInvalidCategoryAmountKind)
	if got != httperr.InternalMessage {
		t.Fatalf("got %q want internal message", got)
	}
}

func TestUserFacingStoreMessage_nilError_mapsToInternal(t *testing.T) {
	t.Parallel()
	got := userFacingStoreMessage(nil)
	if got != httperr.InternalMessage {
		t.Fatalf("got %q want internal message", got)
	}
}

func TestStoreUserMessages_tableIntegrity(t *testing.T) {
	t.Parallel()
	seen := make(map[error]struct{})
	for i, row := range storeUserMessages {
		if row.msg == "" {
			t.Fatalf("row %d: empty message", i)
		}
		if row.sentinel == nil {
			t.Fatalf("row %d: nil sentinel", i)
		}
		if _, dup := seen[row.sentinel]; dup {
			t.Fatalf("row %d: duplicate sentinel %v", i, row.sentinel)
		}
		seen[row.sentinel] = struct{}{}
	}
}
