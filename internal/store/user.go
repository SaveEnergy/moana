package store

import (
	"database/sql"
	"errors"
	"time"
)

// User is an account that can sign in.
type User struct {
	ID            int64
	Email         string
	PasswordHash  []byte
	Role          string
	CreatedAt     time.Time
	HouseholdID   int64
	FirstName     string
	LastName      string
	HouseholdRole string
}

// UserSummary is a user row without secret fields (for admin lists).
type UserSummary struct {
	ID        int64
	Email     string
	Role      string
	CreatedAt time.Time
}

// HouseholdMember is a household row for display (no secrets).
type HouseholdMember struct {
	ID            int64
	Email         string
	FirstName     string
	LastName      string
	HouseholdRole string
}

func scanUser(row *sql.Row) (*User, error) {
	var u User
	var created string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &created,
		&u.HouseholdID, &u.FirstName, &u.LastName, &u.HouseholdRole)
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
	u.CreatedAt = t.UTC()
	return &u, nil
}
