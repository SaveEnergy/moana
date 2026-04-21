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
	if f.FromUTC == nil && f.ToUTC == nil && search == "" && !f.OldestFirst && limit > 0 {
		switch kind {
		case "":
			rows, err := s.DB.QueryContext(ctx, sqlTransactionListRecentDescLimit, householdID, limit)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			return scanTransactionRows(rows, limit)
		case "income":
			rows, err := s.DB.QueryContext(ctx, sqlTransactionListRecentDescLimitIncome, householdID, limit)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			return scanTransactionRows(rows, limit)
		case "expense":
			rows, err := s.DB.QueryContext(ctx, sqlTransactionListRecentDescLimitExpense, householdID, limit)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			return scanTransactionRows(rows, limit)
		}
	}
	var b strings.Builder
	b.Grow(512)
	b.WriteString(sqlTransactionListFromHousehold)
	// At most: household, from, to, 2×search, limit.
	args := make([]any, 0, 6)
	args = append(args, householdID)
	args = appendOccurredAtRange(&b, args, f.FromUTC, f.ToUTC)
	// Only "income" / "expense" add a sign filter; unknown kind strings are ignored (same as Kind "").
	if frag := sqlAmountKindFilter(kind); frag != "" {
		b.WriteString(frag)
	}
	if search != "" {
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
