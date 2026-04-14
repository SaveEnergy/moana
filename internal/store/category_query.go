package store

import (
	"context"
	"database/sql"
	"errors"
)

// ListCategories returns all categories for a user, ordered by name.
func (s *Store) ListCategories(ctx context.Context, userID int64) ([]Category, error) {
	rows, err := s.DB.QueryContext(ctx, `
SELECT id, user_id, name, IFNULL(icon, ''), IFNULL(color, '') FROM categories WHERE user_id = ? ORDER BY name COLLATE NOCASE`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Icon, &c.Color); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// GetCategoryByID returns a category owned by the user, or nil.
func (s *Store) GetCategoryByID(ctx context.Context, userID, categoryID int64) (*Category, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT id, user_id, name, IFNULL(icon, ''), IFNULL(color, '') FROM categories WHERE id = ? AND user_id = ?`, categoryID, userID)
	var c Category
	err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.Icon, &c.Color)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
