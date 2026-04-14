package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

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
