// Package timeutil provides calendar/month ranges, trailing local-day windows, date parsing,
// and SQLite TEXT timestamp encoding (RFC3339Nano UTC in sqlite.go) shared by the store, dashboard,
// history, and htmlview. Calendar boundaries live in calendar.go; rolling windows in trailing_days.go;
// YYYY-MM-DD parsing and [LoadLocation] in range.go (see range_test.go); local calendar date keys in local_date_key.go.
// [OrUTC] maps nil *Location to UTC.
package timeutil
