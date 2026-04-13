package timeutil

import (
	"time"
)

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

// DayRangeUTCFromLocalDates parses YYYY-MM-DD in loc and returns inclusive UTC range for those local calendar days.
func DayRangeUTCFromLocalDates(loc *time.Location, fromDate, toDate string) (fromUTC, toUTC time.Time, err error) {
	from, err := time.ParseInLocation("2006-01-02", fromDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err := time.ParseInLocation("2006-01-02", toDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if to.Before(from) {
		from, to = to, from
	}
	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
	endDay := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, loc)
	end := endDay.Add(24*time.Hour - time.Nanosecond)
	return start.UTC(), end.UTC(), nil
}

// LoadLocation returns UTC if name is empty or invalid.
func LoadLocation(name string) *time.Location {
	if name == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
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
