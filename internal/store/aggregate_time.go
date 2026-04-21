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

// appendOccurredAtRange appends AND t.occurred_at bounds when from/to are non-nil.
func appendOccurredAtRange(b *strings.Builder, args []any, fromUTC, toUTC *time.Time) []any {
	if fromUTC != nil {
		b.WriteString(sqlFilterOccurredAtFrom)
		args = append(args, timeutil.FormatSQLiteUTC(*fromUTC))
	}
	if toUTC != nil {
		b.WriteString(sqlFilterOccurredAtTo)
		args = append(args, timeutil.FormatSQLiteUTC(*toUTC))
	}
	return args
}
