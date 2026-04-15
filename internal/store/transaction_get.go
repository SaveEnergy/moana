package store

import (
	"context"
	"database/sql"
	"errors"
)

// GetTransactionByID returns a transaction visible to the household (same household as the row owner), or nil if not found.
func (s *Store) GetTransactionByID(ctx context.Context, householdID, id int64) (*Transaction, error) {
	row := s.DB.QueryRowContext(ctx, sqlTransactionSelectFromHousehold+`
WHERE t.id = ? AND owner.household_id = ?`, id, householdID)
	var t Transaction
	var occ, cre string
	var catID sql.NullInt64
	err := row.Scan(&t.ID, &t.UserID, &t.AmountCents, &occ, &t.Description, &catID, &t.CategoryName, &t.CategoryIcon, &cre)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := hydrateTransaction(&t, occ, cre, catID); err != nil {
		return nil, err
	}
	return &t, nil
}
