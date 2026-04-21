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

// ParseSQLiteTimestampUTC parses a SQLite TEXT timestamp and returns the same instant in UTC.
// Use for model fields that are stored and compared in UTC ([FormatSQLiteUTC]).
func ParseSQLiteTimestampUTC(s string) (time.Time, error) {
	t, err := ParseSQLiteTimestamp(s)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

// FormatSQLiteUTC formats t as UTC RFC3339Nano for SQLite TEXT columns and query parameters.
func FormatSQLiteUTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// NowSQLiteUTC is FormatSQLiteUTC(time.Now()).
func NowSQLiteUTC() string {
	return FormatSQLiteUTC(time.Now())
}
