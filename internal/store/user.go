package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// User is an account that can sign in.
type User struct {
	ID           int64
	Email        string
	PasswordHash []byte
	Role         string
	Timezone     string
	CreatedAt    time.Time
}

// GetUserByEmail returns a user by email (case-insensitive).
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT id, email, password_hash, role, timezone, created_at
FROM users WHERE email = ? COLLATE NOCASE`, email)
	return scanUser(row)
}

// GetUserByID returns a user by id.
func (s *Store) GetUserByID(ctx context.Context, id int64) (*User, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT id, email, password_hash, role, timezone, created_at
FROM users WHERE id = ?`, id)
	return scanUser(row)
}

func scanUser(row *sql.Row) (*User, error) {
	var u User
	var created string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Timezone, &created)
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

// CreateUser inserts a new user (CLI).
func (s *Store) CreateUser(ctx context.Context, email string, passwordHash []byte, role, timezone string) (int64, error) {
	if timezone == "" {
		timezone = "UTC"
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO users (email, password_hash, role, timezone, created_at)
VALUES (?, ?, ?, ?, ?)`, email, passwordHash, role, timezone, now)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// UpdateUserPassword sets a new password hash.
func (s *Store) UpdateUserPassword(ctx context.Context, userID int64, passwordHash []byte) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, userID)
	return err
}

// UserSummary is a user row without secret fields (for admin lists).
type UserSummary struct {
	ID        int64
	Email     string
	Role      string
	Timezone  string
	CreatedAt time.Time
}

// ListUsers returns all accounts ordered by id (admin).
func (s *Store) ListUsers(ctx context.Context) ([]UserSummary, error) {
	rows, err := s.DB.QueryContext(ctx, `
SELECT id, email, role, timezone, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserSummary
	for rows.Next() {
		var u UserSummary
		var created string
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.Timezone, &created); err != nil {
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
