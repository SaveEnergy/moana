package store

import (
	"context"

	"moana/internal/timeutil"
)

// ListUsers returns all accounts ordered by id (admin).
func (s *Store) ListUsers(ctx context.Context) ([]UserSummary, error) {
	rows, err := s.DB.QueryContext(ctx, `
SELECT id, email, role, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserSummary
	for rows.Next() {
		var u UserSummary
		var created string
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &created); err != nil {
			return nil, err
		}
		t, err := timeutil.ParseSQLiteTimestamp(created)
		if err != nil {
			return nil, err
		}
		u.CreatedAt = t.UTC()
		out = append(out, u)
	}
	return out, rows.Err()
}

// ListHouseholdMembers returns users in the same household.
func (s *Store) ListHouseholdMembers(ctx context.Context, householdID int64) ([]HouseholdMember, error) {
	rows, err := s.DB.QueryContext(ctx, `
SELECT id, email, IFNULL(first_name, ''), IFNULL(last_name, ''), household_role
FROM users WHERE household_id = ? ORDER BY id`, householdID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []HouseholdMember
	for rows.Next() {
		var m HouseholdMember
		if err := rows.Scan(&m.ID, &m.Email, &m.FirstName, &m.LastName, &m.HouseholdRole); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
