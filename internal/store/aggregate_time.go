package store

import (
	"strings"
	"time"

	"moana/internal/timeutil"
)

// SQL AND fragments for optional [t.occurred_at] half-open bounds (used by listing, aggregates, movement).
const (
	sqlFilterOccurredAtFrom = ` AND t.occurred_at >= ?`
	sqlFilterOccurredAtTo   = ` AND t.occurred_at <= ?`
)

func appendOccurredAtRangeSQL(b *strings.Builder, fromUTC, toUTC *time.Time) {
	if fromUTC != nil {
		b.WriteString(sqlFilterOccurredAtFrom)
	}
	if toUTC != nil {
		b.WriteString(sqlFilterOccurredAtTo)
	}
}

func appendOccurredAtRangeArgs(args []any, fromUTC, toUTC *time.Time) []any {
	if fromUTC != nil {
		args = append(args, timeutil.FormatSQLiteUTC(*fromUTC))
	}
	if toUTC != nil {
		args = append(args, timeutil.FormatSQLiteUTC(*toUTC))
	}
	return args
}

// appendOccurredAtRange appends AND t.occurred_at bounds when from/to are non-nil.
func appendOccurredAtRange(b *strings.Builder, args []any, fromUTC, toUTC *time.Time) []any {
	appendOccurredAtRangeSQL(b, fromUTC, toUTC)
	return appendOccurredAtRangeArgs(args, fromUTC, toUTC)
}

func sqlWithOccurredAtRange(prefix string, fromUTC, toUTC *time.Time, suffix string) string {
	var b strings.Builder
	b.Grow(len(prefix) + len(sqlFilterOccurredAtFrom) + len(sqlFilterOccurredAtTo) + len(suffix))
	b.WriteString(prefix)
	appendOccurredAtRangeSQL(&b, fromUTC, toUTC)
	b.WriteString(suffix)
	return b.String()
}
