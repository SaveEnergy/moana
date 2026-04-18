package testutil

import (
	"testing"
	"time"
)

func TestUTCDateString(t *testing.T) {
	t.Parallel()
	tm := time.Date(2026, 3, 5, 18, 0, 0, 0, time.FixedZone("east", 3*3600))
	if got, want := UTCDateString(tm), "2026-03-05"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
