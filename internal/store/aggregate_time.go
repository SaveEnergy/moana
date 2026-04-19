package store

import (
	"strings"
	"time"

	"moana/internal/timeutil"
)

// appendOccurredAtRange appends AND t.occurred_at bounds when from/to are non-nil.
func appendOccurredAtRange(b *strings.Builder, args []any, fromUTC, toUTC *time.Time) []any {
	if fromUTC != nil {
		b.WriteString(` AND t.occurred_at >= ?`)
		args = append(args, timeutil.FormatSQLiteUTC(*fromUTC))
	}
	if toUTC != nil {
		b.WriteString(` AND t.occurred_at <= ?`)
		args = append(args, timeutil.FormatSQLiteUTC(*toUTC))
	}
	return args
}
