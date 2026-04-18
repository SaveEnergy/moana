package store

import (
	"testing"
	"time"
)

func TestApproxLocalDayMapCap_symmetricOrder(t *testing.T) {
	t.Parallel()
	a := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	if got1, got2 := approxLocalDayMapCap(a, b), approxLocalDayMapCap(b, a); got1 != got2 {
		t.Fatalf("got %d vs %d", got1, got2)
	}
}

func TestApproxLocalDayMapCap_clamped(t *testing.T) {
	t.Parallel()
	a := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	b := a.Add(10 * 365 * 24 * time.Hour)
	if n := approxLocalDayMapCap(a, b); n != 500 {
		t.Fatalf("want cap 500, got %d", n)
	}
	if n := approxLocalDayMapCap(a, a); n < 8 {
		t.Fatalf("want min 8, got %d", n)
	}
}
