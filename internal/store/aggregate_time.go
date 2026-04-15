package store

import (
	"time"

	"moana/internal/timeutil"
)

// appendOccurredAtRange appends AND t.occurred_at bounds when from/to are non-nil.
func appendOccurredAtRange(q string, args []any, fromUTC, toUTC *time.Time) (string, []any) {
	if fromUTC != nil {
		q += ` AND t.occurred_at >= ?`
		args = append(args, timeutil.FormatSQLiteUTC(*fromUTC))
	}
	if toUTC != nil {
		q += ` AND t.occurred_at <= ?`
		args = append(args, timeutil.FormatSQLiteUTC(*toUTC))
	}
	return q, args
}
