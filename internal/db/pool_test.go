package db

import "testing"

func TestMaxOpenConns_allowsConcurrentReaders(t *testing.T) {
	t.Parallel()
	if MaxOpenConns < 2 {
		t.Fatalf("MaxOpenConns=%d: pool must allow >1 connection for overlapping WAL reads; keep [Open] wired to this constant", MaxOpenConns)
	}
}
