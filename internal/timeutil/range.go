package timeutil

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// Successful [time.LoadLocation] results are cached by name. The stdlib returns a new
// *[time.Location] on each call, so without this cache every HTTP request would reload
// zoneinfo for the same cookie value (e.g. Europe/Berlin). Concurrent cache misses for
// the same zone name share one [time.LoadLocation] via [singleflight.Group].
var (
	loadLocationCache sync.Map // string -> *time.Location
	loadLocationGroup singleflight.Group
)

// DayRangeUTCFromLocalDates parses YYYY-MM-DD in loc and returns inclusive UTC range for those local calendar days.
// A nil loc is treated as UTC ([time.ParseInLocation] panics if loc is nil).
func DayRangeUTCFromLocalDates(loc *time.Location, fromDate, toDate string) (fromUTC, toUTC time.Time, err error) {
	loc = OrUTC(loc)
	from, err := time.ParseInLocation(time.DateOnly, fromDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err := time.ParseInLocation(time.DateOnly, toDate, loc)
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
	if v, ok := loadLocationCache.Load(name); ok {
		return v.(*time.Location)
	}
	v, _, _ := loadLocationGroup.Do(name, func() (interface{}, error) {
		if v, ok := loadLocationCache.Load(name); ok {
			return v.(*time.Location), nil
		}
		loc, err := time.LoadLocation(name)
		if err != nil {
			return time.UTC, nil
		}
		loadLocationCache.Store(name, loc)
		return loc, nil
	})
	return v.(*time.Location)
}
