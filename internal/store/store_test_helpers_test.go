package store

import (
	"context"
	"errors"
	"testing"

	"moana/internal/db"
)

// testStore opens an in-memory SQLite DB, applies migrations, wraps it in a [Store], and registers t.Cleanup.
func testStore(t *testing.T) *Store {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return New(database)
}

// alreadyCancelledContext returns a [context.Context] that is cancelled before use.
// Pair with [assertErrIsContextCanceled] to verify database/sql propagates cancellation.
func alreadyCancelledContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func assertErrIsContextCanceled(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error when context already cancelled")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v want %v", err, context.Canceled)
	}
}
