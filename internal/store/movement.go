package store

import (
	"context"
	"math"
	"strings"
	"time"

	"moana/internal/money"
	"moana/internal/timeutil"
)

// DailyAbsMovementByLocalDate returns total absolute cents moved per calendar day in loc (sum of |amount_cents| per day) for the household.
// Per-day totals use overflow-checked addition; if the sum exceeds int64, the bucket saturates at [math.MaxInt64].
func (s *Store) DailyAbsMovementByLocalDate(ctx context.Context, householdID int64, fromUTC, toUTC time.Time, loc *time.Location) (map[string]int64, error) {
	loc = timeutil.OrUTC(loc)
	var b strings.Builder
	b.Grow(256)
	b.WriteString(`SELECT t.occurred_at, t.amount_cents `)
	b.WriteString(sqlFromHouseholdTx)
	args := make([]any, 0, 3)
	args = append(args, householdID)
	args = appendOccurredAtRange(&b, args, &fromUTC, &toUTC)
	rows, err := s.DB.QueryContext(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// At most one entry per local calendar day overlapping [fromUTC, toUTC].
	out := make(map[string]int64, approxLocalDayMapCap(fromUTC, toUTC))
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
		cents = money.AbsCents(cents)
		day := timeutil.LocalCalendarDateKey(t, loc)
		prev := out[day]
		next, ok := money.AddCents(prev, cents)
		if !ok {
			out[day] = math.MaxInt64
		} else {
			out[day] = next
		}
	}
	return out, rows.Err()
}

// approxLocalDayMapCap is a coarse upper bound on distinct local dates for bucketing a UTC interval.
func approxLocalDayMapCap(fromUTC, toUTC time.Time) int {
	if toUTC.Before(fromUTC) {
		fromUTC, toUTC = toUTC, fromUTC
	}
	span := toUTC.Sub(fromUTC)
	n := int(span/(24*time.Hour)) + 3
	if n < 8 {
		n = 8
	}
	if n > 500 {
		n = 500
	}
	return n
}
