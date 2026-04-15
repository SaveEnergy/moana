package store

import "context"

// GetUserByEmail returns a user by email (case-insensitive).
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.DB.QueryRowContext(ctx, sqlUserSelectFull+` WHERE email = ? COLLATE NOCASE`, email)
	return scanUser(row)
}

// GetUserByID returns a user by id.
func (s *Store) GetUserByID(ctx context.Context, id int64) (*User, error) {
	row := s.DB.QueryRowContext(ctx, sqlUserSelectFull+` WHERE id = ?`, id)
	return scanUser(row)
}
