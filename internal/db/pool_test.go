package db

import (
	"testing"
	"time"
)

func TestMaxIdleConnLifetime_reasonable(t *testing.T) {
	t.Parallel()
	if maxIdleConnLifetime <= 0 {
		t.Fatal("maxIdleConnLifetime must be positive for SetConnMaxIdleTime")
	}
	if maxIdleConnLifetime > 24*time.Hour {
		t.Fatal("sanity: idle timeout unexpectedly large; fix constant or test")
	}
}

func TestMaxOpenConns_allowsConcurrentReaders(t *testing.T) {
	t.Parallel()
	if MaxOpenConns < 2 {
		t.Fatalf("MaxOpenConns=%d: pool must allow >1 connection for overlapping WAL reads; keep [Open] wired to this constant", MaxOpenConns)
	}
}
