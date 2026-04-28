package store

import (
	"context"
	"database/sql"
	"strings"
)

// ListTransactions returns transactions for the household, ordered by occurred_at (and id).
func (s *Store) ListTransactions(ctx context.Context, householdID int64, f TransactionFilter) ([]Transaction, error) {
	limit := f.Limit
	if limit < 0 {
		limit = 0
	}
	kind := strings.TrimSpace(f.Kind)
	search := strings.TrimSpace(f.Search)

	plan := newTransactionListQueryPlan(f.FromUTC, f.ToUTC, kind, search, f.OldestFirst, limit)
	rows, err := s.DB.QueryContext(ctx, plan.SQL(), transactionListQueryArgs(householdID, f.FromUTC, f.ToUTC, search, limit)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTransactionRows(rows, limit)
}

func scanTransactionRows(rows *sql.Rows, limit int) ([]Transaction, error) {
	cap := 64
	if limit > 0 {
		cap = limit
	}
	out := make([]Transaction, 0, cap)
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
