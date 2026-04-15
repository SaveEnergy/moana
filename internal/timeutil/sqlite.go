package timeutil

import "time"

// ParseSQLiteTimestamp parses timestamps stored in SQLite TEXT columns (RFC3339Nano preferred, RFC3339 fallback).
func ParseSQLiteTimestamp(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Parse(time.RFC3339, s)
	}
	return t, nil
}

// FormatSQLiteUTC formats t as UTC RFC3339Nano for SQLite TEXT columns and query parameters.
func FormatSQLiteUTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// NowSQLiteUTC is FormatSQLiteUTC(time.Now()).
func NowSQLiteUTC() string {
	return FormatSQLiteUTC(time.Now())
}
