package store

import (
	"context"
	"database/sql"

	"moana/internal/timeutil"
)

// CreateUser inserts a new user with their own household (CLI / bootstrap).
func (s *Store) CreateUser(ctx context.Context, email string, passwordHash []byte, role string) (int64, error) {
	email = normalizeUserEmail(email)
	if email == "" {
		return 0, ErrInvalidUserEmail
	}
	now := timeutil.NowSQLiteUTC()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, sqlHouseholdInsertDefaultName, now)
	if err != nil {
		return 0, err
	}
	hid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	res2, err := tx.ExecContext(ctx, sqlUserInsert, email, passwordHash, role, now, hid, "owner")
	if err != nil {
		if sqliteUniqueError(err) {
			return 0, ErrDuplicateUserEmail
		}
		return 0, err
	}
	id, err := res2.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateUserPassword sets a new password hash.
func (s *Store) UpdateUserPassword(ctx context.Context, userID int64, passwordHash []byte) error {
	res, err := s.DB.ExecContext(ctx, sqlUserUpdatePassword, passwordHash, userID)
	if err != nil {
		return err
	}
	return execExactlyOneRow(res, sql.ErrNoRows)
}

// UpdateUserProfile updates name for the signed-in user.
func (s *Store) UpdateUserProfile(ctx context.Context, userID int64, firstName, lastName string) error {
	res, err := s.DB.ExecContext(ctx, sqlUserUpdateProfile, firstName, lastName, userID)
	if err != nil {
		return err
	}
	return execExactlyOneRow(res, sql.ErrNoRows)
}
