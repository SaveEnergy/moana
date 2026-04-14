package store

import (
	"context"
	"time"
)

// DailyAbsMovementByLocalDate returns total absolute cents moved per calendar day in loc (sum of |amount_cents| per day).
func (s *Store) DailyAbsMovementByLocalDate(ctx context.Context, userID int64, fromUTC, toUTC time.Time, loc *time.Location) (map[string]int64, error) {
	q := `SELECT occurred_at, amount_cents FROM transactions WHERE user_id = ? AND occurred_at >= ? AND occurred_at <= ?`
	args := []any{userID, fromUTC.UTC().Format(time.RFC3339Nano), toUTC.UTC().Format(time.RFC3339Nano)}
	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]int64)
	for rows.Next() {
		var occ string
		var cents int64
		if err := rows.Scan(&occ, &cents); err != nil {
			return nil, err
		}
		t, err := parseTime(occ)
		if err != nil {
			return nil, err
		}
		if cents < 0 {
			cents = -cents
		}
		day := t.In(loc).Format("2006-01-02")
		out[day] += cents
	}
	return out, rows.Err()
}
