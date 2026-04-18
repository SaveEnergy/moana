package handlers

import (
	"time"

	"moana/internal/timeutil"
)

// formatLocalCalendarDate returns YYYY-MM-DD for t in loc (nil → UTC; see [timeutil.LocalCalendarDateKey]).
func formatLocalCalendarDate(t time.Time, loc *time.Location) string {
	return timeutil.LocalCalendarDateKey(t, loc)
}

// todayLocalCalendarDate returns today's calendar date in loc (nil → UTC).
func todayLocalCalendarDate(loc *time.Location) string {
	return formatLocalCalendarDate(time.Now(), loc)
}
