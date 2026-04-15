package store

import (
	"context"
	"errors"
)

// CreateCategory adds a category for the household.
func (s *Store) CreateCategory(ctx context.Context, householdID int64, name, icon, color string) (int64, error) {
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO categories (household_id, name, icon, color) VALUES (?, ?, ?, ?)`, householdID, name, icon, color)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateCategory sets name, icon, and color for a category in the household.
func (s *Store) UpdateCategory(ctx context.Context, householdID, categoryID int64, name, icon, color string) error {
	res, err := s.DB.ExecContext(ctx, `
UPDATE categories SET name = ?, icon = ?, color = ? WHERE id = ? AND household_id = ?`, name, icon, color, categoryID, householdID)
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

// DeleteCategory removes a category if it belongs to the household.
func (s *Store) DeleteCategory(ctx context.Context, householdID, categoryID int64) error {
	res, err := s.DB.ExecContext(ctx, `
DELETE FROM categories WHERE id = ? AND household_id = ?`, categoryID, householdID)
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
