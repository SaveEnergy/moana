package store

import (
	"context"
	"database/sql"
	"errors"
)

// ListCategories returns all categories for a household, ordered by name.
func (s *Store) ListCategories(ctx context.Context, householdID int64) ([]Category, error) {
	rows, err := s.DB.QueryContext(ctx,
		sqlCategorySelectFull+` WHERE household_id = ? ORDER BY name COLLATE NOCASE`, householdID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Category, 0, 32)
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.HouseholdID, &c.Name, &c.Icon, &c.Color); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// GetCategoryByID returns a category in the household, or nil.
func (s *Store) GetCategoryByID(ctx context.Context, householdID, categoryID int64) (*Category, error) {
	row := s.DB.QueryRowContext(ctx,
		sqlCategorySelectFull+` WHERE id = ? AND household_id = ?`, categoryID, householdID)
	var c Category
	err := row.Scan(&c.ID, &c.HouseholdID, &c.Name, &c.Icon, &c.Color)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
