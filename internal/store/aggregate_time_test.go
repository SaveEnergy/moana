package store

import (
	"strings"
	"testing"
	"time"
)

func TestAppendOccurredAtRange(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("neither", func(t *testing.T) {
		t.Parallel()
		q, args := appendOccurredAtRange("SELECT 1", []any{42}, nil, nil)
		if q != "SELECT 1" {
			t.Fatalf("q %q", q)
		}
		if len(args) != 1 || args[0] != 42 {
			t.Fatalf("args %v", args)
		}
	})

	t.Run("from_only", func(t *testing.T) {
		t.Parallel()
		q, args := appendOccurredAtRange("SELECT 1", []any{1}, &from, nil)
		if !strings.Contains(q, "t.occurred_at >=") || strings.Contains(q, "t.occurred_at <=") {
			t.Fatalf("q %q", q)
		}
		if len(args) != 2 {
			t.Fatalf("args %v", args)
		}
	})

	t.Run("to_only", func(t *testing.T) {
		t.Parallel()
		q, args := appendOccurredAtRange("SELECT 1", []any{1}, nil, &to)
		if !strings.Contains(q, "t.occurred_at <=") || strings.Contains(q, "t.occurred_at >=") {
			t.Fatalf("q %q", q)
		}
		if len(args) != 2 {
			t.Fatalf("args %v", args)
		}
	})

	t.Run("both", func(t *testing.T) {
		t.Parallel()
		q, args := appendOccurredAtRange("SELECT 1", []any{1}, &from, &to)
		if !strings.Contains(q, "t.occurred_at >=") || !strings.Contains(q, "t.occurred_at <=") {
			t.Fatalf("q %q", q)
		}
		if len(args) != 3 {
			t.Fatalf("args %v", args)
		}
	})
}
