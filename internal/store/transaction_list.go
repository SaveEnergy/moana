package store

import (
	"context"
	"database/sql"

	"moana/internal/timeutil"
)

// ListTransactions returns transactions for the household, ordered by occurred_at (and id).
func (s *Store) ListTransactions(ctx context.Context, householdID int64, f TransactionFilter) ([]Transaction, error) {
	q := sqlTransactionSelectFromHousehold + `
WHERE owner.household_id = ?`
	args := []any{householdID}
	if f.FromUTC != nil {
		q += ` AND t.occurred_at >= ?`
		args = append(args, timeutil.FormatSQLiteUTC(*f.FromUTC))
	}
	if f.ToUTC != nil {
		q += ` AND t.occurred_at <= ?`
		args = append(args, timeutil.FormatSQLiteUTC(*f.ToUTC))
	}
	switch f.Kind {
	case "income":
		q += ` AND t.amount_cents > 0`
	case "expense":
		q += ` AND t.amount_cents < 0`
	}
	if f.Search != "" {
		term := "%" + f.Search + "%"
		q += ` AND (t.description LIKE ? OR COALESCE(c.name, '') LIKE ?)`
		args = append(args, term, term)
	}
	if f.OldestFirst {
		q += ` ORDER BY t.occurred_at ASC, t.id ASC`
	} else {
		q += ` ORDER BY t.occurred_at DESC, t.id DESC`
	}
	if f.Limit > 0 {
		q += ` LIMIT ?`
		args = append(args, f.Limit)
	}

	rows, err := s.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Transaction
	for rows.Next() {
		var t Transaction
		var occ, cre string
		var catID sql.NullInt64
		if err := rows.Scan(&t.ID, &t.UserID, &t.AmountCents, &occ, &t.Description, &catID, &t.CategoryName, &t.CategoryIcon, &cre); err != nil {
			return nil, err
		}
		if err := hydrateTransaction(&t, occ, cre, catID); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
