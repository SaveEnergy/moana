package timeutil

import "time"

// OrUTC returns loc if non-nil; otherwise [time.UTC] (several time APIs panic on nil *Location).
func OrUTC(loc *time.Location) *time.Location {
	if loc == nil {
		return time.UTC
	}
	return loc
}
