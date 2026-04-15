// Package timeutil provides calendar/month ranges, trailing local-day windows, date parsing,
// and SQLite TEXT timestamp encoding (RFC3339Nano UTC in sqlite.go) shared by the store, dashboard,
// history, and htmlview. Calendar boundaries live in calendar.go; rolling windows in trailing_days.go;
// YYYY-MM-DD parsing in range.go.
package timeutil
