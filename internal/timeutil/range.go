package timeutil

import "time"

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
