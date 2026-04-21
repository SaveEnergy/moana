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
	var b strings.Builder
	b.Grow(512)
	b.WriteString(sqlTransactionSelectFromHousehold)
	b.WriteString("\nWHERE owner.household_id = ?")
	// At most: household, from, to, 2×search, limit.
	args := make([]any, 0, 6)
	args = append(args, householdID)
	args = appendOccurredAtRange(&b, args, f.FromUTC, f.ToUTC)
	// Only "income" / "expense" add a sign filter; unknown kind strings are ignored (same as Kind "").
	kind := strings.TrimSpace(f.Kind)
	if frag := sqlAmountKindFilter(kind); frag != "" {
		b.WriteString(frag)
	}
	if search := strings.TrimSpace(f.Search); search != "" {
		term := "%" + escapeSQLLikePattern(search) + "%"
		b.WriteString(` AND (t.description LIKE ? ESCAPE '!' OR COALESCE(c.name, '') LIKE ? ESCAPE '!')`)
		args = append(args, term, term)
	}
	if f.OldestFirst {
		b.WriteString(` ORDER BY t.occurred_at ASC, t.id ASC`)
	} else {
		b.WriteString(` ORDER BY t.occurred_at DESC, t.id DESC`)
	}
	if limit > 0 {
		b.WriteString(` LIMIT ?`)
		args = append(args, limit)
	}

	rows, err := s.DB.QueryContext(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Transaction
	switch {
	case limit > 0:
		out = make([]Transaction, 0, limit)
	default:
		out = make([]Transaction, 0, 64)
	}
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
