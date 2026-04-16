package timeutil

import (
	"testing"
	"time"
)

func TestSQLiteTimestampRoundTrip(t *testing.T) {
	t.Parallel()
	orig := time.Date(2024, 3, 15, 14, 30, 0, 123456789, time.UTC)
	s := FormatSQLiteUTC(orig)
	got, err := ParseSQLiteTimestamp(s)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(orig) {
		t.Fatalf("got %v want %v", got, orig)
	}
}

func TestParseSQLiteTimestampRFC3339Fallback(t *testing.T) {
	t.Parallel()
	s := "2024-03-15T14:30:00Z"
	got, err := ParseSQLiteTimestamp(s)
	if err != nil {
		t.Fatal(err)
	}
	if got.Year() != 2024 || got.Month() != time.March || got.Day() != 15 {
		t.Fatalf("got %v", got)
	}
}

func TestParseSQLiteTimestamp_invalid(t *testing.T) {
	t.Parallel()
	_, err := ParseSQLiteTimestamp("totally not a date")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNowSQLiteUTC_roundTrips(t *testing.T) {
	t.Parallel()
	s := NowSQLiteUTC()
	got, err := ParseSQLiteTimestamp(s)
	if err != nil {
		t.Fatalf("NowSQLiteUTC() = %q: %v", s, err)
	}
	if got.IsZero() {
		t.Fatal("parsed time is zero")
	}
}
