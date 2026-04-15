package store

import (
	"context"

	"moana/internal/timeutil"
)

// UpdateHouseholdName sets the display name for a household.
func (s *Store) UpdateHouseholdName(ctx context.Context, householdID int64, name string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE households SET name = ? WHERE id = ?`, name, householdID)
	return err
}

// CreateHouseholdMember adds a user with role member to an existing household.
func (s *Store) CreateHouseholdMember(ctx context.Context, householdID int64, email string, passwordHash []byte) (int64, error) {
	now := timeutil.NowSQLiteUTC()
	res, err := s.DB.ExecContext(ctx, sqlUserInsert, email, passwordHash, "user", now, householdID, "member")
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// DetachUserToSoloHousehold moves a user into a new household as owner (remove from shared household).
func (s *Store) DetachUserToSoloHousehold(ctx context.Context, userID int64) error {
	now := timeutil.NowSQLiteUTC()
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
