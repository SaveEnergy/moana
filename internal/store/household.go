package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Household is a shared account group.
type Household struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

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
	t, err := time.Parse(time.RFC3339Nano, created)
	if err != nil {
		t, err = time.Parse(time.RFC3339, created)
	}
	if err != nil {
		return nil, err
	}
	h.CreatedAt = t.UTC()
	return &h, nil
}

// UpdateHouseholdName sets the display name for a household.
func (s *Store) UpdateHouseholdName(ctx context.Context, householdID int64, name string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE households SET name = ? WHERE id = ?`, name, householdID)
	return err
}

// CountHouseholdMembers returns how many users belong to a household.
func (s *Store) CountHouseholdMembers(ctx context.Context, householdID int64) (int64, error) {
	var n int64
	err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE household_id = ?`, householdID).Scan(&n)
	return n, err
}
