package store

import (
	"context"
	"database/sql"
	"time"
)

// SumAmountCents returns the sum of amount_cents for the household in the optional time range.
func (s *Store) SumAmountCents(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time) (int64, error) {
	return s.SumAmountCentsByKind(ctx, householdID, fromUTC, toUTC, "")
}

// SumAmountCentsByKind sums amounts in [from, to]; kind is "", "income", or "expense".
func (s *Store) SumAmountCentsByKind(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time, kind string) (int64, error) {
	q := `SELECT COALESCE(SUM(t.amount_cents), 0) FROM transactions t
INNER JOIN users u ON u.id = t.user_id
WHERE u.household_id = ?`
	args := []any{householdID}
	q, args = appendOccurredAtRange(q, args, fromUTC, toUTC)
	switch kind {
	case "income":
		q += ` AND t.amount_cents > 0`
	case "expense":
		q += ` AND t.amount_cents < 0`
	}
	var sum sql.NullInt64
	err := s.DB.QueryRowContext(ctx, q, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}
	if !sum.Valid {
		return 0, nil
	}
	return sum.Int64, nil
}
