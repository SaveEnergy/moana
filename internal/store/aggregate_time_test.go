package store

import (
	"strings"
	"testing"
	"time"

	"moana/internal/timeutil"
)

func TestAppendOccurredAtRange(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("neither", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{42}, nil, nil)
		if b.String() != "SELECT 1" {
			t.Fatalf("q %q", b.String())
		}
		if len(args) != 1 || args[0] != 42 {
			t.Fatalf("args %v", args)
		}
	})

	t.Run("from_only", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{1}, &from, nil)
		want := "SELECT 1" + sqlFilterOccurredAtFrom
		if b.String() != want {
			t.Fatalf("q %q want %q", b.String(), want)
		}
		if len(args) != 2 {
			t.Fatalf("args %v", args)
		}
		if args[1] != timeutil.FormatSQLiteUTC(from) {
			t.Fatalf("bound arg %v", args[1])
		}
	})

	t.Run("to_only", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{1}, nil, &to)
		want := "SELECT 1" + sqlFilterOccurredAtTo
		if b.String() != want {
			t.Fatalf("q %q want %q", b.String(), want)
		}
		if len(args) != 2 {
			t.Fatalf("args %v", args)
		}
		if args[1] != timeutil.FormatSQLiteUTC(to) {
			t.Fatalf("bound arg %v", args[1])
		}
	})

	t.Run("both", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{1}, &from, &to)
		want := "SELECT 1" + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo
		if b.String() != want {
			t.Fatalf("q %q want %q", b.String(), want)
		}
		if len(args) != 3 {
			t.Fatalf("args %v", args)
		}
		if args[1] != timeutil.FormatSQLiteUTC(from) || args[2] != timeutil.FormatSQLiteUTC(to) {
			t.Fatalf("bounds %v %v", args[1], args[2])
		}
	})
}
