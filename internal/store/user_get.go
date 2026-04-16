package store

import "context"

// GetUserByEmail returns a user by email (case-insensitive). Surrounding whitespace is ignored.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	email = normalizeUserEmail(email)
	if email == "" {
		return nil, nil
	}
	row := s.DB.QueryRowContext(ctx, sqlUserSelectFull+` WHERE email = ? COLLATE NOCASE`, email)
	return scanUser(row)
}

// GetUserByID returns a user by id.
func (s *Store) GetUserByID(ctx context.Context, id int64) (*User, error) {
	row := s.DB.QueryRowContext(ctx, sqlUserSelectFull+` WHERE id = ?`, id)
	return scanUser(row)
}
