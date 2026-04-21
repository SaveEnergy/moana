package store

import (
	"context"
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
	if fromUTC == nil && toUTC == nil {
		err = s.DB.QueryRowContext(ctx, sqlSumIncomeExpenseBase, householdID).Scan(&incomeSum, &expenseSum)
		return incomeSum, expenseSum, err
	}
	switch {
	case fromUTC != nil && toUTC != nil:
		err = s.DB.QueryRowContext(ctx, sqlSumIncomeExpenseInRangeBoth,
			householdID,
			timeutil.FormatSQLiteUTC(*fromUTC),
			timeutil.FormatSQLiteUTC(*toUTC),
		).Scan(&incomeSum, &expenseSum)
	case fromUTC != nil:
		err = s.DB.QueryRowContext(ctx, sqlSumIncomeExpenseInRangeFromOnly,
			householdID,
			timeutil.FormatSQLiteUTC(*fromUTC),
		).Scan(&incomeSum, &expenseSum)
	default:
		err = s.DB.QueryRowContext(ctx, sqlSumIncomeExpenseInRangeToOnly,
			householdID,
			timeutil.FormatSQLiteUTC(*toUTC),
		).Scan(&incomeSum, &expenseSum)
	}
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
	aF, aT := timeutil.FormatSQLiteUTC(aFrom), timeutil.FormatSQLiteUTC(aTo)
	bF, bT := timeutil.FormatSQLiteUTC(bFrom), timeutil.FormatSQLiteUTC(bTo)
	args := []any{aF, aT, aF, aT, bF, bT, bF, bT, householdID}
	err = s.DB.QueryRowContext(ctx, sqlSumRunningTotalAndIncomeExpenseInTwoRanges, args...).Scan(&running, &aIncome, &aExpense, &bIncome, &bExpense)
	return running, aIncome, aExpense, bIncome, bExpense, err
}

// SumAmountCentsByKind sums amounts in [from, to]; kind is "", "income", or "expense" (ASCII trim applied, same as listing filters).
func (s *Store) SumAmountCentsByKind(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time, kind string) (int64, error) {
	if fromUTC == nil && toUTC == nil {
		var sum int64
		switch frag := sqlAmountKindFilter(kind); frag {
		case "":
			err := s.DB.QueryRowContext(ctx, sqlSumAmountAllHousehold, householdID).Scan(&sum)
			return sum, err
		case sqlFilterAmountIncome:
			err := s.DB.QueryRowContext(ctx, sqlSumAmountIncomeHouseholdNoBounds, householdID).Scan(&sum)
			return sum, err
		case sqlFilterAmountExpense:
			err := s.DB.QueryRowContext(ctx, sqlSumAmountExpenseHouseholdNoBounds, householdID).Scan(&sum)
			return sum, err
		}
	}
	frag := sqlAmountKindFilter(kind)
	var sum int64
	var err error
	switch {
	case fromUTC != nil && toUTC != nil:
		f, t := timeutil.FormatSQLiteUTC(*fromUTC), timeutil.FormatSQLiteUTC(*toUTC)
		switch frag {
		case "":
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeBoth, householdID, f, t).Scan(&sum)
		case sqlFilterAmountIncome:
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeBothIncome, householdID, f, t).Scan(&sum)
		case sqlFilterAmountExpense:
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeBothExpense, householdID, f, t).Scan(&sum)
		}
	case fromUTC != nil:
		f := timeutil.FormatSQLiteUTC(*fromUTC)
		switch frag {
		case "":
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeFromOnly, householdID, f).Scan(&sum)
		case sqlFilterAmountIncome:
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeFromOnlyIncome, householdID, f).Scan(&sum)
		case sqlFilterAmountExpense:
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeFromOnlyExpense, householdID, f).Scan(&sum)
		}
	default:
		t := timeutil.FormatSQLiteUTC(*toUTC)
		switch frag {
		case "":
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeToOnly, householdID, t).Scan(&sum)
		case sqlFilterAmountIncome:
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeToOnlyIncome, householdID, t).Scan(&sum)
		case sqlFilterAmountExpense:
			err = s.DB.QueryRowContext(ctx, sqlSumAmountHouseholdRangeToOnlyExpense, householdID, t).Scan(&sum)
		}
	}
	if err != nil {
		return 0, err
	}
	return sum, nil
}
