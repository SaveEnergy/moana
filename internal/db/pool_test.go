package db

import "testing"

func TestSQLiteMaxOpenConns_allowsConcurrentReaders(t *testing.T) {
	t.Parallel()
	if sqliteMaxOpenConns < 2 {
		t.Fatalf("sqliteMaxOpenConns=%d: pool must allow >1 connection for overlapping WAL reads; keep [Open] wired to this constant", sqliteMaxOpenConns)
	}
}
