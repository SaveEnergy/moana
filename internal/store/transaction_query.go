package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// GetTransactionByID returns a transaction owned by userID, or nil if not found.
func (s *Store) GetTransactionByID(ctx context.Context, userID, id int64) (*Transaction, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT t.id, t.user_id, t.amount_cents, t.occurred_at, t.description, t.category_id, COALESCE(c.name, ''), COALESCE(IFNULL(c.icon, ''), ''), t.created_at
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE t.id = ? AND t.user_id = ?`, id, userID)
	var t Transaction
	var occ, cre string
	var catID sql.NullInt64
	err := row.Scan(&t.ID, &t.UserID, &t.AmountCents, &occ, &t.Description, &catID, &t.CategoryName, &t.CategoryIcon, &cre)
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

// ListTransactions returns transactions for the user, ordered by occurred_at (and id).
func (s *Store) ListTransactions(ctx context.Context, userID int64, f TransactionFilter) ([]Transaction, error) {
	q := `
SELECT t.id, t.user_id, t.amount_cents, t.occurred_at, t.description, t.category_id, COALESCE(c.name, ''), COALESCE(IFNULL(c.icon, ''), ''), t.created_at
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
		if err := rows.Scan(&t.ID, &t.UserID, &t.AmountCents, &occ, &t.Description, &catID, &t.CategoryName, &t.CategoryIcon, &cre); err != nil {
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
