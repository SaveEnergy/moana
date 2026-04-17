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

func TestTrailingLocalDaysInclusiveRangeUTC_nonPositiveDaysUsesOne(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	ref := time.Date(2026, 4, 13, 12, 0, 0, 0, loc)
	curS, curE := TrailingLocalDaysInclusiveRangeUTC(loc, ref, 0)
	if curS.In(loc).Format("2006-01-02") != "2026-04-13" || curE.In(loc).Format("2006-01-02") != "2026-04-13" {
		t.Fatalf("current single-day window: %v .. %v", curS.In(loc), curE.In(loc))
	}
	prevS, prevE := PriorTrailingLocalDaysInclusiveRangeUTC(loc, ref, -3)
	if prevS.In(loc).Format("2006-01-02") != "2026-04-12" || prevE.In(loc).Format("2006-01-02") != "2026-04-12" {
		t.Fatalf("prior single-day window: %v .. %v", prevS.In(loc), prevE.In(loc))
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

func TestDayRangeUTCFromLocalDates_singleDayInclusive(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	from, to, err := DayRangeUTCFromLocalDates(loc, "2026-06-01", "2026-06-01")
	if err != nil {
		t.Fatal(err)
	}
	if from != time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("start of day: %v", from)
	}
	if !to.After(from) {
		t.Fatalf("end of inclusive day must be after start: from=%v to=%v", from, to)
	}
}

func TestDayRangeUTCFromLocalDates_swappedArgsMatchForwardRange(t *testing.T) {
	t.Parallel()
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	fwdFrom, fwdTo, err := DayRangeUTCFromLocalDates(loc, "2026-01-10", "2026-01-12")
	if err != nil {
		t.Fatal(err)
	}
	revFrom, revTo, err := DayRangeUTCFromLocalDates(loc, "2026-01-12", "2026-01-10")
	if err != nil {
		t.Fatal(err)
	}
	if !fwdFrom.Equal(revFrom) || !fwdTo.Equal(revTo) {
		t.Fatalf("swapped args should match forward range: fwd %v..%v rev %v..%v", fwdFrom, fwdTo, revFrom, revTo)
	}
}

func TestDayRangeUTCFromLocalDates_nilLocationUsesUTC(t *testing.T) {
	t.Parallel()
	fromNil, toNil, err := DayRangeUTCFromLocalDates(nil, "2026-06-01", "2026-06-01")
	if err != nil {
		t.Fatal(err)
	}
	fromUTC, toUTC, err := DayRangeUTCFromLocalDates(time.UTC, "2026-06-01", "2026-06-01")
	if err != nil {
		t.Fatal(err)
	}
	if !fromNil.Equal(fromUTC) || !toNil.Equal(toUTC) {
		t.Fatalf("nil loc: %v..%v want same as UTC %v..%v", fromNil, toNil, fromUTC, toUTC)
	}
}

func TestTrailingLocalDaysInclusiveRangeUTC_nilLocationMatchesUTC(t *testing.T) {
	t.Parallel()
	ref := time.Date(2026, 4, 13, 15, 0, 0, 0, time.UTC)
	s, e := TrailingLocalDaysInclusiveRangeUTC(nil, ref, 30)
	s2, e2 := TrailingLocalDaysInclusiveRangeUTC(time.UTC, ref, 30)
	if !s.Equal(s2) || !e.Equal(e2) {
		t.Fatalf("nil vs UTC: %v..%v vs %v..%v", s, e, s2, e2)
	}
}

func TestPriorTrailingLocalDaysInclusiveRangeUTC_nilLocationMatchesUTC(t *testing.T) {
	t.Parallel()
	ref := time.Date(2026, 4, 13, 15, 0, 0, 0, time.UTC)
	s, e := PriorTrailingLocalDaysInclusiveRangeUTC(nil, ref, 30)
	s2, e2 := PriorTrailingLocalDaysInclusiveRangeUTC(time.UTC, ref, 30)
	if !s.Equal(s2) || !e.Equal(e2) {
		t.Fatalf("nil vs UTC: %v..%v vs %v..%v", s, e, s2, e2)
	}
}
