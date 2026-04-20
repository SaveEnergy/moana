package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"moana/internal/timeutil"
)

// ErrInvalidCategory is returned when a category id does not exist in the household.
var ErrInvalidCategory = errors.New("invalid category")

// ErrUserNotInHousehold is returned when userID does not belong to householdID (CreateTransaction only).
var ErrUserNotInHousehold = errors.New("user not in household")

// CreateTransaction inserts a transaction. occurredAt must be UTC.
// userID must belong to householdID; the insert is conditional on that (no orphan rows for mismatched ids).
func (s *Store) CreateTransaction(ctx context.Context, userID, householdID int64, amountCents int64, occurredAt time.Time, description string, categoryID *int64) (int64, error) {
	if err := s.validateCategoryOwnership(ctx, householdID, categoryID); err != nil {
		return 0, err
	}
	occ := timeutil.FormatSQLiteUTC(occurredAt)
	now := timeutil.NowSQLiteUTC()
	cat := sqlNullCategoryID(categoryID)
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO transactions (user_id, amount_cents, occurred_at, description, category_id, created_at)
SELECT ?, ?, ?, ?, ?, ?
WHERE EXISTS (SELECT 1 FROM users WHERE id = ? AND household_id = ?)`,
		userID, amountCents, occ, description, cat, now, userID, householdID)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, ErrUserNotInHousehold
	}
	return res.LastInsertId()
}

// UpdateTransaction updates a row in the caller's household. categoryID is validated against the household's categories. occurredAt must be UTC.
// actorUserID must belong to householdID (same enforcement pattern as [Store.CreateTransaction]).
func (s *Store) UpdateTransaction(ctx context.Context, householdID, actorUserID, id int64, amountCents int64, occurredAt time.Time, description string, categoryID *int64) error {
	if err := s.validateCategoryOwnership(ctx, householdID, categoryID); err != nil {
		return err
	}
	occ := timeutil.FormatSQLiteUTC(occurredAt)
	cat := sqlNullCategoryID(categoryID)
	res, err := s.DB.ExecContext(ctx, `
UPDATE transactions SET amount_cents = ?, occurred_at = ?, description = ?, category_id = ?
WHERE id = ? AND user_id IN (SELECT id FROM users WHERE household_id = ?)
AND EXISTS (SELECT 1 FROM users WHERE id = ? AND household_id = ?)`,
		amountCents, occ, description, cat, id, householdID, actorUserID, householdID)
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

func (s *Store) validateCategoryOwnership(ctx context.Context, householdID int64, categoryID *int64) error {
	if categoryID == nil {
		return nil
	}
	cat, err := s.GetCategoryByID(ctx, householdID, *categoryID)
	if err != nil {
		return err
	}
	if cat == nil {
		return ErrInvalidCategory
	}
	return nil
}

func sqlNullCategoryID(categoryID *int64) sql.NullInt64 {
	if categoryID == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *categoryID, Valid: true}
}
