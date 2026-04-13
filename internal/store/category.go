package store

import (
	"context"
	"database/sql"
	"errors"
)

// Category is a user-owned label for transactions.
type Category struct {
	ID     int64
	UserID int64
	Name   string
}

// ListCategories returns all categories for a user, ordered by name.
func (s *Store) ListCategories(ctx context.Context, userID int64) ([]Category, error) {
	rows, err := s.DB.QueryContext(ctx, `
SELECT id, user_id, name FROM categories WHERE user_id = ? ORDER BY name COLLATE NOCASE`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// CreateCategory adds a category for the user.
func (s *Store) CreateCategory(ctx context.Context, userID int64, name string) (int64, error) {
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO categories (user_id, name) VALUES (?, ?)`, userID, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
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

// GetCategoryByID returns a category owned by the user, or nil.
func (s *Store) GetCategoryByID(ctx context.Context, userID, categoryID int64) (*Category, error) {
	row := s.DB.QueryRowContext(ctx, `
SELECT id, user_id, name FROM categories WHERE id = ? AND user_id = ?`, categoryID, userID)
	var c Category
	err := row.Scan(&c.ID, &c.UserID, &c.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
