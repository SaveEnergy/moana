package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Transaction is a single income (or adjustment) entry.
type Transaction struct {
	ID           int64
	UserID       int64
	AmountCents  int64
	OccurredAt   time.Time
	Description  string
	CategoryID   sql.NullInt64
	CategoryName string
	CreatedAt    time.Time
}

// CreateTransaction inserts a transaction. occurredAt must be UTC.
func (s *Store) CreateTransaction(ctx context.Context, userID int64, amountCents int64, occurredAt time.Time, description string, categoryID *int64) (int64, error) {
	if categoryID != nil {
		cat, err := s.GetCategoryByID(ctx, userID, *categoryID)
		if err != nil {
			return 0, err
		}
		if cat == nil {
			return 0, errors.New("invalid category")
		}
	}
	occ := occurredAt.UTC().Format(time.RFC3339Nano)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var cat sql.NullInt64
	if categoryID != nil {
		cat = sql.NullInt64{Int64: *categoryID, Valid: true}
	}
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO transactions (user_id, amount_cents, occurred_at, description, category_id, created_at)
VALUES (?, ?, ?, ?, ?, ?)`, userID, amountCents, occ, description, cat, now)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetTransactionByID returns a transaction owned by userID, or nil if not found.
func (s *Store) GetTransactionByID(ctx context.Context, userID, id int64) (*Transaction, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT t.id, t.user_id, t.amount_cents, t.occurred_at, t.description, t.category_id, COALESCE(c.name, ''), t.created_at
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE t.id = ? AND t.user_id = ?`, id, userID)
	var t Transaction
	var occ, cre string
	var catID sql.NullInt64
	err := row.Scan(&t.ID, &t.UserID, &t.AmountCents, &occ, &t.Description, &catID, &t.CategoryName, &cre)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t.OccurredAt, err = parseTime(occ)
	if err != nil {
		return nil, err
	}
	t.CreatedAt, err = parseTime(cre)
	if err != nil {
		return nil, err
	}
	t.CategoryID = catID
	return &t, nil
}

// UpdateTransaction updates a row owned by userID. occurredAt must be UTC.
func (s *Store) UpdateTransaction(ctx context.Context, userID, id int64, amountCents int64, occurredAt time.Time, description string, categoryID *int64) error {
	if categoryID != nil {
		cat, err := s.GetCategoryByID(ctx, userID, *categoryID)
		if err != nil {
			return err
		}
		if cat == nil {
			return errors.New("invalid category")
		}
	}
	occ := occurredAt.UTC().Format(time.RFC3339Nano)
	var cat sql.NullInt64
	if categoryID != nil {
		cat = sql.NullInt64{Int64: *categoryID, Valid: true}
	}
	res, err := s.DB.ExecContext(ctx, `
UPDATE transactions SET amount_cents = ?, occurred_at = ?, description = ?, category_id = ?
WHERE id = ? AND user_id = ?`, amountCents, occ, description, cat, id, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// TransactionFilter restricts listed transactions.
type TransactionFilter struct {
	FromUTC     *time.Time
	ToUTC       *time.Time
	Kind        string // "", "all", "income", "expense"
	Search      string // matches description or category name (substring)
	OldestFirst bool   // default false = newest first
	Limit       int    // if > 0, cap rows (newest/oldest order preserved)
}

// ListTransactions returns transactions for the user, ordered by occurred_at (and id).
func (s *Store) ListTransactions(ctx context.Context, userID int64, f TransactionFilter) ([]Transaction, error) {
	q := `
SELECT t.id, t.user_id, t.amount_cents, t.occurred_at, t.description, t.category_id, COALESCE(c.name, ''), t.created_at
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE t.user_id = ?`
	args := []any{userID}
	if f.FromUTC != nil {
		q += ` AND t.occurred_at >= ?`
		args = append(args, f.FromUTC.UTC().Format(time.RFC3339Nano))
	}
	if f.ToUTC != nil {
		q += ` AND t.occurred_at <= ?`
		args = append(args, f.ToUTC.UTC().Format(time.RFC3339Nano))
	}
	switch f.Kind {
	case "income":
		q += ` AND t.amount_cents > 0`
	case "expense":
		q += ` AND t.amount_cents < 0`
	}
	if f.Search != "" {
		term := "%" + f.Search + "%"
		q += ` AND (t.description LIKE ? OR COALESCE(c.name, '') LIKE ?)`
		args = append(args, term, term)
	}
	if f.OldestFirst {
		q += ` ORDER BY t.occurred_at ASC, t.id ASC`
	} else {
		q += ` ORDER BY t.occurred_at DESC, t.id DESC`
	}
	if f.Limit > 0 {
		q += ` LIMIT ?`
		args = append(args, f.Limit)
	}

	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Transaction
	for rows.Next() {
		var t Transaction
		var occ, cre string
		var catID sql.NullInt64
		if err := rows.Scan(&t.ID, &t.UserID, &t.AmountCents, &occ, &t.Description, &catID, &t.CategoryName, &cre); err != nil {
			return nil, err
		}
		t.OccurredAt, err = parseTime(occ)
		if err != nil {
			return nil, err
		}
		t.CreatedAt, err = parseTime(cre)
		if err != nil {
			return nil, err
		}
		t.CategoryID = catID
		out = append(out, t)
	}
	return out, rows.Err()
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Parse(time.RFC3339, s)
	}
	return t, nil
}

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
	AmountCents int64 // always positive
}

// ListCategoryAmountsInRange returns per-category totals for income (amount_cents > 0)
// or expense (amount_cents < 0, returned as positive magnitudes), ordered by size.
func (s *Store) ListCategoryAmountsInRange(ctx context.Context, userID int64, fromUTC, toUTC *time.Time, kind string) ([]CategoryAmount, error) {
	if kind != "income" && kind != "expense" {
		return nil, errors.New("kind must be income or expense")
	}
	q := `
SELECT t.category_id, COALESCE(c.name, 'Uncategorized'), COALESCE(SUM(t.amount_cents), 0)
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
	q += ` GROUP BY t.category_id, c.name`
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
		if err := rows.Scan(&ca.CategoryID, &ca.Name, &sum); err != nil {
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
