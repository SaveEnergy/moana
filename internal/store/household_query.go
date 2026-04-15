package store

import (
	"context"
	"database/sql"
	"errors"

	"moana/internal/timeutil"
)

// GetHousehold returns a household by id.
func (s *Store) GetHousehold(ctx context.Context, id int64) (*Household, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT id, name, created_at FROM households WHERE id = ?`, id)
	var h Household
	var created string
	err := row.Scan(&h.ID, &h.Name, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t, err := timeutil.ParseSQLiteTimestamp(created)
	if err != nil {
		return nil, err
	}
	h.CreatedAt = t.UTC()
	return &h, nil
}

// CountHouseholdMembers returns how many users belong to a household.
func (s *Store) CountHouseholdMembers(ctx context.Context, householdID int64) (int64, error) {
	var n int64
	err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE household_id = ?`, householdID).Scan(&n)
	return n, err
}
