package store

import (
	"context"
	"errors"
)

// CreateCategory adds a category for the user.
func (s *Store) CreateCategory(ctx context.Context, userID int64, name, icon, color string) (int64, error) {
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO categories (user_id, name, icon, color) VALUES (?, ?, ?, ?)`, userID, name, icon, color)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateCategory sets name, icon, and color for a category owned by the user.
func (s *Store) UpdateCategory(ctx context.Context, userID, categoryID int64, name, icon, color string) error {
	res, err := s.DB.ExecContext(ctx, `
UPDATE categories SET name = ?, icon = ?, color = ? WHERE id = ? AND user_id = ?`, name, icon, color, categoryID, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("category not found")
	}
	return nil
}

// DeleteCategory removes a category if it belongs to the user.
func (s *Store) DeleteCategory(ctx context.Context, userID, categoryID int64) error {
	res, err := s.DB.ExecContext(ctx, `
DELETE FROM categories WHERE id = ? AND user_id = ?`, categoryID, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("category not found")
	}
	return nil
}
