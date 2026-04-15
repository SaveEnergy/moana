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

// SumIncomeExpenseCentsInRange returns income (positive amounts) and expense (negative amounts) sums
// for the household in one query. Net equals incomeSum + expenseSum.
func (s *Store) SumIncomeExpenseCentsInRange(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time) (incomeSum int64, expenseSum int64, err error) {
	q := `SELECT COALESCE(SUM(CASE WHEN t.amount_cents > 0 THEN t.amount_cents ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN t.amount_cents < 0 THEN t.amount_cents ELSE 0 END), 0) ` + sqlFromHouseholdTx
	args := []any{householdID}
	q, args = appendOccurredAtRange(q, args, fromUTC, toUTC)
	err = s.DB.QueryRowContext(ctx, q, args...).Scan(&incomeSum, &expenseSum)
	return incomeSum, expenseSum, err
}

// SumAmountCentsByKind sums amounts in [from, to]; kind is "", "income", or "expense".
func (s *Store) SumAmountCentsByKind(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time, kind string) (int64, error) {
	q := `SELECT COALESCE(SUM(t.amount_cents), 0) ` + sqlFromHouseholdTx
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
