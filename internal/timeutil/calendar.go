package timeutil

import "time"

// PreviousCalendarMonthRangeUTC returns [start, end] in UTC for the calendar month
// immediately before the one containing ref (ref is typically "now" in UTC).
func PreviousCalendarMonthRangeUTC(loc *time.Location, ref time.Time) (startUTC, endUTC time.Time) {
	localRef := ref.In(loc)
	startThisMonth := time.Date(localRef.Year(), localRef.Month(), 1, 0, 0, 0, 0, loc)
	startPrev := startThisMonth.AddDate(0, -1, 0)
	endPrev := startThisMonth.Add(-time.Nanosecond)
	return startPrev.UTC(), endPrev.UTC()
}

// PreviousCalendarYearRangeUTC returns [start, end] in UTC for the calendar year
// immediately before the one containing ref.
func PreviousCalendarYearRangeUTC(loc *time.Location, ref time.Time) (startUTC, endUTC time.Time) {
	localRef := ref.In(loc)
	startThisYear := time.Date(localRef.Year(), 1, 1, 0, 0, 0, 0, loc)
	startPrev := startThisYear.AddDate(-1, 0, 0)
	endPrev := startThisYear.Add(-time.Nanosecond)
	return startPrev.UTC(), endPrev.UTC()
}

// CalendarMonthRangeUTC returns [start, end] UTC for the calendar month that is
// monthsBack months before the month containing ref (0 = month containing ref).
func CalendarMonthRangeUTC(loc *time.Location, ref time.Time, monthsBack int) (startUTC, endUTC time.Time) {
	localRef := ref.In(loc)
	start := time.Date(localRef.Year(), localRef.Month(), 1, 0, 0, 0, 0, loc).AddDate(0, -monthsBack, 0)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start.UTC(), end.UTC()
}

// CurrentCalendarYearToDateRangeUTC returns [start of year in loc, end of ref's local day] in UTC.
func CurrentCalendarYearToDateRangeUTC(loc *time.Location, ref time.Time) (startUTC, endUTC time.Time) {
	localRef := ref.In(loc)
	start := time.Date(localRef.Year(), 1, 1, 0, 0, 0, 0, loc)
	end := time.Date(localRef.Year(), localRef.Month(), localRef.Day(), 23, 59, 59, 999999999, loc)
	return start.UTC(), end.UTC()
}
