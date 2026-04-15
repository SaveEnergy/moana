package store

import (
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
