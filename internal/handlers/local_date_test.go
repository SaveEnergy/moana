package handlers

import (
	"testing"
	"time"
)

func TestFormatLocalCalendarDate_nilLocationUsesUTC(t *testing.T) {
	t.Parallel()
	tm := time.Date(2026, 3, 20, 15, 30, 0, 0, time.UTC)
	if got := formatLocalCalendarDate(tm, nil); got != "2026-03-20" {
		t.Fatalf("got %q", got)
	}
}

func TestTodayLocalCalendarDate_nilDoesNotPanic(t *testing.T) {
	t.Parallel()
	s := todayLocalCalendarDate(nil)
	if len(s) != 10 || s[4] != '-' || s[7] != '-' {
		t.Fatalf("got %q", s)
	}
}
