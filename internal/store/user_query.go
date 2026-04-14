package store

import (
	"context"
	"time"
)

// GetUserByEmail returns a user by email (case-insensitive).
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT id, email, password_hash, role, created_at, household_id, first_name, last_name, household_role
FROM users WHERE email = ? COLLATE NOCASE`, email)
	return scanUser(row)
}

// GetUserByID returns a user by id.
func (s *Store) GetUserByID(ctx context.Context, id int64) (*User, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT id, email, password_hash, role, created_at, household_id, first_name, last_name, household_role
FROM users WHERE id = ?`, id)
	return scanUser(row)
}

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
		t, err := time.Parse(time.RFC3339Nano, created)
		if err != nil {
			t, err = time.Parse(time.RFC3339, created)
		}
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
