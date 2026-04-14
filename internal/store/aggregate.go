package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// SumAmountCents returns the sum of amount_cents for the user in the optional time range.
func (s *Store) SumAmountCents(ctx context.Context, userID int64, fromUTC, toUTC *time.Time) (int64, error) {
	return s.SumAmountCentsByKind(ctx, userID, fromUTC, toUTC, "")
}

// SumAmountCentsByKind sums amounts in [from, to]; kind is "", "income", or "expense".
func (s *Store) SumAmountCentsByKind(ctx context.Context, userID int64, fromUTC, toUTC *time.Time, kind string) (int64, error) {
	q := `SELECT COALESCE(SUM(amount_cents), 0) FROM transactions WHERE user_id = ?`
	args := []any{userID}
	if fromUTC != nil {
		q += ` AND occurred_at >= ?`
		args = append(args, fromUTC.UTC().Format(time.RFC3339Nano))
	}
	if toUTC != nil {
		q += ` AND occurred_at <= ?`
		args = append(args, toUTC.UTC().Format(time.RFC3339Nano))
	}
	switch kind {
	case "income":
		q += ` AND amount_cents > 0`
	case "expense":
		q += ` AND amount_cents < 0`
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

// CategoryExpense is total negative amount for one category in a range.
type CategoryExpense struct {
	CategoryID   sql.NullInt64
	CategoryName string
	TotalCents   int64
}

// ListTopExpenseCategories returns categories with the largest expenses (most negative sums), up to limit.
func (s *Store) ListTopExpenseCategories(ctx context.Context, userID int64, fromUTC, toUTC *time.Time, limit int) ([]CategoryExpense, error) {
	if limit < 1 {
		limit = 5
	}
	q := `
SELECT t.category_id, COALESCE(c.name, 'Uncategorized'), COALESCE(SUM(t.amount_cents), 0)
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE t.user_id = ? AND t.amount_cents < 0`
	args := []any{userID}
	if fromUTC != nil {
		q += ` AND t.occurred_at >= ?`
		args = append(args, fromUTC.UTC().Format(time.RFC3339Nano))
	}
	if toUTC != nil {
		q += ` AND t.occurred_at <= ?`
		args = append(args, toUTC.UTC().Format(time.RFC3339Nano))
	}
	q += ` GROUP BY t.category_id, c.name ORDER BY SUM(t.amount_cents) ASC LIMIT ?`
	args = append(args, limit)

	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CategoryExpense
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
func (s *Store) ListCategoryAmountsInRange(ctx context.Context, userID int64, fromUTC, toUTC *time.Time, kind string) ([]CategoryAmount, error) {
	if kind != "income" && kind != "expense" {
		return nil, errors.New("kind must be income or expense")
	}
	q := `
SELECT t.category_id, COALESCE(MAX(c.name), 'Uncategorized'), COALESCE(MAX(IFNULL(c.icon, '')), ''), COALESCE(MAX(IFNULL(c.color, '')), ''), COALESCE(SUM(t.amount_cents), 0)
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE t.user_id = ?`
	args := []any{userID}
	if fromUTC != nil {
		q += ` AND t.occurred_at >= ?`
		args = append(args, fromUTC.UTC().Format(time.RFC3339Nano))
	}
	if toUTC != nil {
		q += ` AND t.occurred_at <= ?`
		args = append(args, toUTC.UTC().Format(time.RFC3339Nano))
	}
	if kind == "income" {
		q += ` AND t.amount_cents > 0`
	} else {
		q += ` AND t.amount_cents < 0`
	}
	q += ` GROUP BY t.category_id`
	if kind == "income" {
		q += ` ORDER BY SUM(t.amount_cents) DESC`
	} else {
		q += ` ORDER BY SUM(t.amount_cents) ASC`
	}

	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CategoryAmount
	for rows.Next() {
		var ca CategoryAmount
		var sum int64
		if err := rows.Scan(&ca.CategoryID, &ca.Name, &ca.Icon, &ca.Color, &sum); err != nil {
			return nil, err
		}
		if kind == "expense" {
			sum = -sum
		}
		ca.AmountCents = sum
		out = append(out, ca)
	}
	return out, rows.Err()
}
