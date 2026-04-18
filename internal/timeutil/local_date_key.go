package timeutil

import "time"

// LocalCalendarDateKey returns the local calendar date for t in loc as YYYY-MM-DD ([time.DateOnly]).
// A nil loc is treated as UTC ([OrUTC]).
func LocalCalendarDateKey(t time.Time, loc *time.Location) string {
	return t.In(OrUTC(loc)).Format(time.DateOnly)
}
