package dbutil

import (
	"testing"

	"moana/internal/store"
)

// MustOpenMemStore opens SQLite :memory: (migrated), wraps a [store.Store], and registers t.Cleanup to close the DB.
// Use from unit tests outside [moana/internal/store] (that package uses its own test helper to avoid an import cycle with dbutil).
func MustOpenMemStore(t *testing.T) *store.Store {
	t.Helper()
	st, sqlDB, err := OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return st
}
