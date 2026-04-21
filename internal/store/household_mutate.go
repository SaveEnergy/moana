package store

import (
	"context"
	"database/sql"
	"errors"

	"moana/internal/timeutil"
)

// ErrDuplicateUserEmail is returned when inserting a user would violate the unique email constraint.
var ErrDuplicateUserEmail = errors.New("duplicate user email")

// ErrInvalidUserEmail is returned when an email is empty or only whitespace after normalization.
var ErrInvalidUserEmail = errors.New("invalid user email")

// UpdateHouseholdName sets the display name for a household.
func (s *Store) UpdateHouseholdName(ctx context.Context, householdID int64, name string) error {
	res, err := s.DB.ExecContext(ctx, sqlHouseholdUpdateName, name, householdID)
	if err != nil {
		return err
	}
	return execExactlyOneRow(res, sql.ErrNoRows)
}

// CreateHouseholdMember adds a user with role member to an existing household.
func (s *Store) CreateHouseholdMember(ctx context.Context, householdID int64, email string, passwordHash []byte) (int64, error) {
	email = normalizeUserEmail(email)
	if email == "" {
		return 0, ErrInvalidUserEmail
	}
	now := timeutil.NowSQLiteUTC()
	res, err := s.DB.ExecContext(ctx, sqlUserInsert, email, passwordHash, "user", now, householdID, "member")
	if err != nil {
		if sqliteUniqueError(err) {
			return 0, ErrDuplicateUserEmail
		}
		return 0, err
	}
	return res.LastInsertId()
}

// DetachUserToSoloHousehold moves a user into a new household as owner (remove from shared household).
// If no user row matches userID, it returns [sql.ErrNoRows] and rolls back the new household insert.
func (s *Store) DetachUserToSoloHousehold(ctx context.Context, userID int64) error {
	now := timeutil.NowSQLiteUTC()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.ExecContext(ctx, sqlHouseholdInsertDefaultName, now)
	if err != nil {
		return err
	}
	hid, err := res.LastInsertId()
	if err != nil {
		return err
	}
	res2, err := tx.ExecContext(ctx, sqlUserDetachToSoloHousehold, hid, userID)
	if err != nil {
		return err
	}
	if err := execExactlyOneRow(res2, sql.ErrNoRows); err != nil {
		return err
	}
	return tx.Commit()
}
