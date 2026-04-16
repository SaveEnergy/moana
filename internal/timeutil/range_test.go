package timeutil

import (
	"testing"
	"time"
)

func TestPreviousCalendarMonthRangeUTC(t *testing.T) {
	t.Parallel()
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	// 15 Feb 2026 local -> previous month is Jan 2026
	ref := time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC)
	start, end := PreviousCalendarMonthRangeUTC(loc, ref)
	localStart := start.In(loc)
	localEnd := end.In(loc)
	if localStart.Month() != time.January || localStart.Day() != 1 || localStart.Hour() != 0 {
		t.Fatalf("start in local: %v", localStart)
	}
	if localEnd.Month() != time.January || localEnd.Day() != 31 {
		t.Fatalf("end in local: %v", localEnd)
	}
	if !start.Before(end) {
		t.Fatalf("start %v not before end %v", start, end)
	}
}

func TestPreviousCalendarYearRangeUTC(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	ref := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	start, end := PreviousCalendarYearRangeUTC(loc, ref)
	if start.Year() != 2025 || start.Month() != time.January || start.Day() != 1 {
		t.Fatalf("start: %v", start)
	}
	if end.Year() != 2025 || end.Month() != time.December || end.Day() != 31 {
		t.Fatalf("end: %v", end)
	}
}

func TestTrailingLocalDaysInclusiveRangeUTC(t *testing.T) {
	t.Parallel()
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	ref := time.Date(2026, 4, 13, 15, 0, 0, 0, loc)
	curS, curE := TrailingLocalDaysInclusiveRangeUTC(loc, ref, 30)
	prevS, prevE := PriorTrailingLocalDaysInclusiveRangeUTC(loc, ref, 30)
	if curS.In(loc).Format("2006-01-02") != "2026-03-15" || curE.In(loc).Format("2006-01-02") != "2026-04-13" {
		t.Fatalf("current window start=%v end=%v", curS.In(loc), curE.In(loc))
	}
	if prevS.In(loc).Format("2006-01-02") != "2026-02-13" || prevE.In(loc).Format("2006-01-02") != "2026-03-14" {
		t.Fatalf("prior window start=%v end=%v", prevS.In(loc), prevE.In(loc))
	}
}

func TestLoadLocation(t *testing.T) {
	t.Parallel()
	if LoadLocation("") != time.UTC {
		t.Fatal("empty name must be UTC")
	}
	if LoadLocation("not-a-real-zone-name-xyz") != time.UTC {
		t.Fatal("invalid name must be UTC")
	}
	if LoadLocation("UTC") != time.UTC {
		t.Fatal(`LoadLocation("UTC") must be time.UTC`)
	}
	berlin := LoadLocation("Europe/Berlin")
	if berlin == nil || berlin == time.UTC {
		t.Fatal("valid zone must not be UTC")
	}
}

func TestDayRangeUTCFromLocalDates(t *testing.T) {
	t.Parallel()
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	from, to, err := DayRangeUTCFromLocalDates(loc, "2026-01-10", "2026-01-12")
	if err != nil {
		t.Fatal(err)
	}
	if from.In(loc).Day() != 10 || to.In(loc).Day() != 12 {
		t.Fatalf("from=%v to=%v", from.In(loc), to.In(loc))
	}
}
