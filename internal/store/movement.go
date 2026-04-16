package store

import (
	"context"
	"time"

	"moana/internal/timeutil"
)

// DailyAbsMovementByLocalDate returns total absolute cents moved per calendar day in loc (sum of |amount_cents| per day) for the household.
func (s *Store) DailyAbsMovementByLocalDate(ctx context.Context, householdID int64, fromUTC, toUTC time.Time, loc *time.Location) (map[string]int64, error) {
	q := `SELECT t.occurred_at, t.amount_cents ` + sqlFromHouseholdTx + ` AND t.occurred_at >= ? AND t.occurred_at <= ?`
	args := make([]any, 0, 3)
	args = append(args, householdID, timeutil.FormatSQLiteUTC(fromUTC), timeutil.FormatSQLiteUTC(toUTC))
	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// At most one entry per calendar day in range (~366 for rolling year heatmap).
	out := make(map[string]int64, 400)
	for rows.Next() {
		var occ string
		var cents int64
		if err := rows.Scan(&occ, &cents); err != nil {
			return nil, err
		}
		t, err := timeutil.ParseSQLiteTimestamp(occ)
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
