package store

import (
	"context"

	"moana/internal/timeutil"
)

// CreateUser inserts a new user with their own household (CLI / bootstrap).
func (s *Store) CreateUser(ctx context.Context, email string, passwordHash []byte, role string) (int64, error) {
	now := timeutil.NowSQLiteUTC()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `INSERT INTO households (name, created_at) VALUES ('My household', ?)`, now)
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
	_, err := s.DB.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, userID)
	return err
}

// UpdateUserProfile updates name for the signed-in user.
func (s *Store) UpdateUserProfile(ctx context.Context, userID int64, firstName, lastName string) error {
	_, err := s.DB.ExecContext(ctx, `
UPDATE users SET first_name = ?, last_name = ? WHERE id = ?`, firstName, lastName, userID)
	return err
}
