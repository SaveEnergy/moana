package store

import (
	"context"
	"database/sql"
	"time"

	"moana/internal/timeutil"
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
	// cap: household + optional from/to.
	args := make([]any, 0, 3)
	args = append(args, householdID)
	q, args = appendOccurredAtRange(q, args, fromUTC, toUTC)
	err = s.DB.QueryRowContext(ctx, q, args...).Scan(&incomeSum, &expenseSum)
	return incomeSum, expenseSum, err
}

// SumIncomeExpenseCentsInTwoRanges returns income and expense sums for two closed intervals in one query.
// Ranges use the same semantics as [Store.SumIncomeExpenseCentsInRange] (inclusive bounds on occurred_at).
// It is implemented via [Store.SumRunningTotalAndIncomeExpenseInTwoRanges] (same scan, discards running total).
func (s *Store) SumIncomeExpenseCentsInTwoRanges(ctx context.Context, householdID int64, aFrom, aTo, bFrom, bTo time.Time) (aIncome, aExpense, bIncome, bExpense int64, err error) {
	_, aIncome, aExpense, bIncome, bExpense, err = s.SumRunningTotalAndIncomeExpenseInTwoRanges(ctx, householdID, aFrom, aTo, bFrom, bTo)
	return aIncome, aExpense, bIncome, bExpense, err
}

// SumRunningTotalAndIncomeExpenseInTwoRanges returns the all-time net sum for the household (same as
// [Store.SumAmountCents] with no date filter) plus income and expense totals for two closed intervals,
// in a single scan. Used by the dashboard to avoid an extra round trip.
func (s *Store) SumRunningTotalAndIncomeExpenseInTwoRanges(ctx context.Context, householdID int64, aFrom, aTo, bFrom, bTo time.Time) (running int64, aIncome, aExpense, bIncome, bExpense int64, err error) {
	q := `SELECT
  COALESCE(SUM(t.amount_cents), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents > 0 THEN t.amount_cents ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents < 0 THEN t.amount_cents ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents > 0 THEN t.amount_cents ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents < 0 THEN t.amount_cents ELSE 0 END), 0)
` + sqlFromHouseholdTx
	aF, aT := timeutil.FormatSQLiteUTC(aFrom), timeutil.FormatSQLiteUTC(aTo)
	bF, bT := timeutil.FormatSQLiteUTC(bFrom), timeutil.FormatSQLiteUTC(bTo)
	args := []any{aF, aT, aF, aT, bF, bT, bF, bT, householdID}
	err = s.DB.QueryRowContext(ctx, q, args...).Scan(&running, &aIncome, &aExpense, &bIncome, &bExpense)
	return running, aIncome, aExpense, bIncome, bExpense, err
}

// SumAmountCentsByKind sums amounts in [from, to]; kind is "", "income", or "expense".
func (s *Store) SumAmountCentsByKind(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time, kind string) (int64, error) {
	q := `SELECT COALESCE(SUM(t.amount_cents), 0) ` + sqlFromHouseholdTx
	// cap: household + optional from/to.
	args := make([]any, 0, 3)
	args = append(args, householdID)
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
