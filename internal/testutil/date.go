package testutil

import "time"

// UTCDateString returns t as a UTC calendar date (YYYY-MM-DD) for HTML date inputs in tests.
func UTCDateString(t time.Time) string {
	return t.UTC().Format(time.DateOnly)
}
