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
	if f.FromUTC == nil && f.ToUTC == nil && search == "" && limit > 0 {
		if q, ok := staticListTransactionsQuery(f.OldestFirst, kind); ok {
			rows, err := s.DB.QueryContext(ctx, q, householdID, limit)
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
		b.WriteString(sqlTransactionListSearchLike)
		args = append(args, term, term)
	}
	if f.OldestFirst {
		b.WriteString(sqlTransactionListOrderAsc)
	} else {
		b.WriteString(sqlTransactionListOrderDesc)
	}
	if limit > 0 {
		b.WriteString(sqlTransactionListLimitSuffix)
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

// staticListTransactionsQuery returns fixed SQL for household + sort + optional income|expense filter + LIMIT.
// kind must be trimmed; unknown kinds use the dynamic query path in [Store.ListTransactions].
func staticListTransactionsQuery(oldestFirst bool, kind string) (query string, ok bool) {
	switch kind {
	case "":
		if oldestFirst {
			return sqlTransactionListOldestAscLimit, true
		}
		return sqlTransactionListRecentDescLimit, true
	case "income":
		if oldestFirst {
			return sqlTransactionListOldestAscLimitIncome, true
		}
		return sqlTransactionListRecentDescLimitIncome, true
	case "expense":
		if oldestFirst {
			return sqlTransactionListOldestAscLimitExpense, true
		}
		return sqlTransactionListRecentDescLimitExpense, true
	default:
		return "", false
	}
}
