package timeutil

import "time"

// TrailingLocalDaysInclusiveRangeUTC returns [start, end] UTC for inclusiveDays local calendar days
// ending on ref's local calendar day (inclusive). E.g. inclusiveDays=30 is the last 30 local days.
func TrailingLocalDaysInclusiveRangeUTC(loc *time.Location, ref time.Time, inclusiveDays int) (startUTC, endUTC time.Time) {
	if inclusiveDays < 1 {
		inclusiveDays = 1
	}
	localRef := ref.In(loc)
	todayLocal := time.Date(localRef.Year(), localRef.Month(), localRef.Day(), 0, 0, 0, 0, loc)
	startDay := todayLocal.AddDate(0, 0, -(inclusiveDays - 1))
	start := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, loc)
	end := time.Date(todayLocal.Year(), todayLocal.Month(), todayLocal.Day(), 23, 59, 59, 999999999, loc)
	return start.UTC(), end.UTC()
}

// PriorTrailingLocalDaysInclusiveRangeUTC returns the same-length window immediately before
// TrailingLocalDaysInclusiveRangeUTC (for period-over-period comparisons).
func PriorTrailingLocalDaysInclusiveRangeUTC(loc *time.Location, ref time.Time, inclusiveDays int) (startUTC, endUTC time.Time) {
	if inclusiveDays < 1 {
		inclusiveDays = 1
	}
	localRef := ref.In(loc)
	todayLocal := time.Date(localRef.Year(), localRef.Month(), localRef.Day(), 0, 0, 0, 0, loc)
	endPrevDay := todayLocal.AddDate(0, 0, -inclusiveDays)
	startPrevDay := todayLocal.AddDate(0, 0, -(2*inclusiveDays - 1))
	start := time.Date(startPrevDay.Year(), startPrevDay.Month(), startPrevDay.Day(), 0, 0, 0, 0, loc)
	end := time.Date(endPrevDay.Year(), endPrevDay.Month(), endPrevDay.Day(), 23, 59, 59, 999999999, loc)
	return start.UTC(), end.UTC()
}
