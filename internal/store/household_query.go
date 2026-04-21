package store

import (
	"context"
	"database/sql"
	"errors"

	"moana/internal/timeutil"
)

// GetHousehold returns a household by id.
func (s *Store) GetHousehold(ctx context.Context, id int64) (*Household, error) {
	row := s.DB.QueryRowContext(ctx, sqlHouseholdGetByID, id)
	var h Household
	var created string
	err := row.Scan(&h.ID, &h.Name, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t, err := timeutil.ParseSQLiteTimestampUTC(created)
	if err != nil {
		return nil, err
	}
	h.CreatedAt = t
	return &h, nil
}

// CountHouseholdMembers returns how many users belong to a household.
func (s *Store) CountHouseholdMembers(ctx context.Context, householdID int64) (int64, error) {
	var n int64
	err := s.DB.QueryRowContext(ctx, sqlHouseholdCountMembers, householdID).Scan(&n)
	return n, err
}
