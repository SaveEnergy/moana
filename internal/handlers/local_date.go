package handlers

import (
	"time"

	"moana/internal/timeutil"
)

// formatLocalCalendarDate returns YYYY-MM-DD for t in loc (nil → UTC; [time.Time.In] panics on nil *Location).
func formatLocalCalendarDate(t time.Time, loc *time.Location) string {
	return t.In(timeutil.OrUTC(loc)).Format("2006-01-02")
}

// todayLocalCalendarDate returns today's calendar date in loc (nil → UTC).
func todayLocalCalendarDate(loc *time.Location) string {
	return formatLocalCalendarDate(time.Now(), loc)
}
