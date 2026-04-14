package store

import (
	"context"
	"time"
)

// CreateUser inserts a new user with their own household (CLI / bootstrap).
func (s *Store) CreateUser(ctx context.Context, email string, passwordHash []byte, role string) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
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
	res2, err := tx.ExecContext(ctx, `
INSERT INTO users (email, password_hash, role, created_at, household_id, first_name, last_name, household_role)
VALUES (?, ?, ?, ?, ?, '', '', 'owner')`, email, passwordHash, role, now, hid)
	if err != nil {
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

// CreateHouseholdMember adds a user with role member to an existing household.
func (s *Store) CreateHouseholdMember(ctx context.Context, householdID int64, email string, passwordHash []byte) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO users (email, password_hash, role, created_at, household_id, first_name, last_name, household_role)
VALUES (?, ?, 'user', ?, ?, '', '', 'member')`, email, passwordHash, now, householdID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// DetachUserToSoloHousehold moves a user into a new household as owner (remove from shared household).
func (s *Store) DetachUserToSoloHousehold(ctx context.Context, userID int64) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.ExecContext(ctx, `INSERT INTO households (name, created_at) VALUES ('My household', ?)`, now)
	if err != nil {
		return err
	}
	hid, err := res.LastInsertId()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE users SET household_id = ?, household_role = 'owner' WHERE id = ?`, hid, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
