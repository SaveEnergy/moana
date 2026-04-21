package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"moana/internal/money"
)

// ErrInvalidCategoryAmountKind is returned when [Store.ListCategoryAmountsInRange] is called with kind other than "income" or "expense".
var ErrInvalidCategoryAmountKind = errors.New("kind must be income or expense")

// CategoryExpense is total negative amount for one category in a range.
type CategoryExpense struct {
	CategoryID   sql.NullInt64
	CategoryName string
	TotalCents   int64
}

// ListTopExpenseCategories returns categories with the largest expenses (most negative sums), up to limit.
func (s *Store) ListTopExpenseCategories(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time, limit int) ([]CategoryExpense, error) {
	if limit < 1 {
		limit = 5
	}
	if fromUTC == nil && toUTC == nil {
		rows, err := s.DB.QueryContext(ctx, sqlListTopExpenseCategoriesNoDate, householdID, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanCategoryExpenseRows(rows, limit)
	}
	var b strings.Builder
	b.Grow(512)
	b.WriteString(sqlListTopExpenseCategoriesPrefix)
	// cap: household + optional from/to + limit.
	args := make([]any, 0, 4)
	args = append(args, householdID)
	args = appendOccurredAtRange(&b, args, fromUTC, toUTC)
	b.WriteString(sqlListTopExpenseCategoriesSuffix)
	args = append(args, limit)

	rows, err := s.DB.QueryContext(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCategoryExpenseRows(rows, limit)
}

func scanCategoryExpenseRows(rows *sql.Rows, limit int) ([]CategoryExpense, error) {
	out := make([]CategoryExpense, 0, limit)
	for rows.Next() {
		var ce CategoryExpense
		if err := rows.Scan(&ce.CategoryID, &ce.CategoryName, &ce.TotalCents); err != nil {
			return nil, err
		}
		out = append(out, ce)
	}
	return out, rows.Err()
}

// CategoryAmount is a positive total for one category in a date range.
type CategoryAmount struct {
	CategoryID  sql.NullInt64
	Name        string
	Icon        string // from categories.icon when joined; "" if uncategorized or unset
	Color       string // hex from categories.color; "" if unset
	AmountCents int64  // always positive
}

// ListCategoryAmountsInRange returns per-category totals for income (amount_cents > 0)
// or expense (amount_cents < 0, returned as positive magnitudes), ordered by size.
func (s *Store) ListCategoryAmountsInRange(ctx context.Context, householdID int64, fromUTC, toUTC *time.Time, kind string) ([]CategoryAmount, error) {
	kind = strings.TrimSpace(kind)
	if kind != "income" && kind != "expense" {
		return nil, ErrInvalidCategoryAmountKind
	}
	if fromUTC == nil && toUTC == nil {
		var q string
		switch kind {
		case "income":
			q = sqlListCategoryAmountsIncomeFullHousehold
		case "expense":
			q = sqlListCategoryAmountsExpenseFullHousehold
		}
		rows, err := s.DB.QueryContext(ctx, q, householdID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanCategoryAmountRows(rows)
	}
	var b strings.Builder
	b.Grow(768)
	b.WriteString(sqlListCategoryAmountsSelectPrefix)
	// cap: household + optional from/to.
	args := make([]any, 0, 3)
	args = append(args, householdID)
	args = appendOccurredAtRange(&b, args, fromUTC, toUTC)
	b.WriteString(sqlAmountKindFilter(kind))
	b.WriteString(sqlListCategoryAmountsGroupBy)
	if kind == "income" {
		b.WriteString(sqlListCategoryAmountsOrderIncome)
	} else {
		b.WriteString(sqlListCategoryAmountsOrderExpense)
	}

	rows, err := s.DB.QueryContext(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCategoryAmountRows(rows)
}

func scanCategoryAmountRows(rows *sql.Rows) ([]CategoryAmount, error) {
	out := make([]CategoryAmount, 0, 32)
	for rows.Next() {
		var ca CategoryAmount
		var sum int64
		if err := rows.Scan(&ca.CategoryID, &ca.Name, &ca.Icon, &ca.Color, &sum); err != nil {
			return nil, err
		}
		// Expense sums are negative in SQL; use magnitude (plain negation overflows at MinInt64).
		ca.AmountCents = money.AbsCents(sum)
		out = append(out, ca)
	}
	return out, rows.Err()
}
