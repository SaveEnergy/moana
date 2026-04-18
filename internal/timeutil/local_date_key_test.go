package timeutil

import (
	"testing"
	"time"
)

func TestLocalCalendarDateKey_BerlinCrossesUTCDate(t *testing.T) {
	t.Parallel()
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	// 2026-01-15 23:30 UTC → local 2026-01-16 CET.
	tm := time.Date(2026, 1, 15, 23, 30, 0, 0, time.UTC)
	if got, want := LocalCalendarDateKey(tm, loc), "2026-01-16"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestLocalCalendarDateKey_nilLocationUsesUTC(t *testing.T) {
	t.Parallel()
	tm := time.Date(2026, 6, 10, 15, 0, 0, 0, time.UTC)
	if got, want := LocalCalendarDateKey(tm, nil), "2026-06-10"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
