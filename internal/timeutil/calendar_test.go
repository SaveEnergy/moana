package timeutil

import (
	"testing"
	"time"
)

func TestCalendarMonthRangeUTC_monthsBackZeroIsMonthOfRef(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	ref := time.Date(2026, 3, 15, 12, 0, 0, 0, loc)
	start, end := CalendarMonthRangeUTC(loc, ref, 0)
	if got := start.In(loc).Format("2006-01-02"); got != "2026-03-01" {
		t.Fatalf("start %v (%s)", start, got)
	}
	if end.In(loc).Month() != time.March || end.In(loc).Day() != 31 {
		t.Fatalf("end should be last moment of March, got %v", end.In(loc))
	}
	if !start.Before(end) {
		t.Fatalf("start %v not before end %v", start, end)
	}
}

func TestCalendarMonthRangeUTC_monthsBackOneIsPriorMonth(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	ref := time.Date(2026, 3, 15, 0, 0, 0, 0, loc)
	start, end := CalendarMonthRangeUTC(loc, ref, 1)
	if got := start.In(loc).Format("2006-01-02"); got != "2026-02-01" {
		t.Fatalf("start %s", got)
	}
	if end.In(loc).Month() != time.February {
		t.Fatalf("end month %v", end.In(loc).Month())
	}
}

func TestCurrentCalendarYearToDateRangeUTC(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	ref := time.Date(2026, 3, 15, 14, 30, 45, 0, loc)
	start, end := CurrentCalendarYearToDateRangeUTC(loc, ref)
	if start.In(loc).Format("2006-01-02") != "2026-01-01" {
		t.Fatalf("start %v", start.In(loc))
	}
	if end.In(loc).Format("2006-01-02") != "2026-03-15" {
		t.Fatalf("end local day %v", end.In(loc))
	}
	if end.Nanosecond() != 999999999 {
		t.Fatalf("end nanosecond %d", end.Nanosecond())
	}
}
