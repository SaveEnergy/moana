package store

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"moana/internal/timeutil"
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
	if search == "" && (f.FromUTC != nil || f.ToUTC != nil) {
		if q, ok := staticListTransactionsDatedNoSearchQuery(f.FromUTC, f.ToUTC, f.OldestFirst, kind, limit); ok {
			args := listTransactionDatedNoSearchArgs(householdID, f.FromUTC, f.ToUTC, limit)
			rows, err := s.DB.QueryContext(ctx, q, args...)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			return scanTransactionRows(rows, limit)
		}
	}
	if f.FromUTC == nil && f.ToUTC == nil && search != "" {
		if q, ok := staticListTransactionsSearchNoDateQuery(kind, f.OldestFirst, limit); ok {
			term := "%" + escapeSQLLikePattern(search) + "%"
			args := []any{householdID, term, term}
			if limit > 0 {
				args = append(args, limit)
			}
			rows, err := s.DB.QueryContext(ctx, q, args...)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			return scanTransactionRows(rows, limit)
		}
	}
	if search != "" && f.FromUTC != nil && f.ToUTC != nil {
		if q, ok := staticListTransactionsDatedBothSearchQuery(kind, f.OldestFirst, limit); ok {
			term := "%" + escapeSQLLikePattern(search) + "%"
			args := []any{householdID, timeutil.FormatSQLiteUTC(*f.FromUTC), timeutil.FormatSQLiteUTC(*f.ToUTC), term, term}
			if limit > 0 {
				args = append(args, limit)
			}
			rows, err := s.DB.QueryContext(ctx, q, args...)
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

func listTransactionDatedNoSearchArgs(householdID int64, fromUTC, toUTC *time.Time, limit int) []any {
	switch {
	case fromUTC != nil && toUTC != nil:
		a := []any{householdID, timeutil.FormatSQLiteUTC(*fromUTC), timeutil.FormatSQLiteUTC(*toUTC)}
		if limit > 0 {
			a = append(a, limit)
		}
		return a
	case fromUTC != nil:
		a := []any{householdID, timeutil.FormatSQLiteUTC(*fromUTC)}
		if limit > 0 {
			a = append(a, limit)
		}
		return a
	default:
		a := []any{householdID, timeutil.FormatSQLiteUTC(*toUTC)}
		if limit > 0 {
			a = append(a, limit)
		}
		return a
	}
}

// staticListTransactionsDatedNoSearchQuery returns fixed SQL when [Store.ListTransactions] has date bounds and no search (optional LIMIT when limit > 0).
func staticListTransactionsDatedNoSearchQuery(fromUTC, toUTC *time.Time, oldestFirst bool, kind string, limit int) (query string, ok bool) {
	withLimit := limit > 0
	frag := sqlAmountKindFilter(kind)
	both := fromUTC != nil && toUTC != nil
	fromOnly := fromUTC != nil && toUTC == nil
	switch {
	case both:
		switch frag {
		case "":
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedBothNoKindAscLimit, true
				}
				return sqlTransactionListDatedBothNoKindAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedBothNoKindDescLimit, true
			}
			return sqlTransactionListDatedBothNoKindDesc, true
		case sqlFilterAmountIncome:
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedBothIncomeAscLimit, true
				}
				return sqlTransactionListDatedBothIncomeAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedBothIncomeDescLimit, true
			}
			return sqlTransactionListDatedBothIncomeDesc, true
		case sqlFilterAmountExpense:
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedBothExpenseAscLimit, true
				}
				return sqlTransactionListDatedBothExpenseAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedBothExpenseDescLimit, true
			}
			return sqlTransactionListDatedBothExpenseDesc, true
		}
	case fromOnly:
		switch frag {
		case "":
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedFromOnlyNoKindAscLimit, true
				}
				return sqlTransactionListDatedFromOnlyNoKindAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedFromOnlyNoKindDescLimit, true
			}
			return sqlTransactionListDatedFromOnlyNoKindDesc, true
		case sqlFilterAmountIncome:
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedFromOnlyIncomeAscLimit, true
				}
				return sqlTransactionListDatedFromOnlyIncomeAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedFromOnlyIncomeDescLimit, true
			}
			return sqlTransactionListDatedFromOnlyIncomeDesc, true
		case sqlFilterAmountExpense:
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedFromOnlyExpenseAscLimit, true
				}
				return sqlTransactionListDatedFromOnlyExpenseAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedFromOnlyExpenseDescLimit, true
			}
			return sqlTransactionListDatedFromOnlyExpenseDesc, true
		}
	default:
		switch frag {
		case "":
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedToOnlyNoKindAscLimit, true
				}
				return sqlTransactionListDatedToOnlyNoKindAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedToOnlyNoKindDescLimit, true
			}
			return sqlTransactionListDatedToOnlyNoKindDesc, true
		case sqlFilterAmountIncome:
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedToOnlyIncomeAscLimit, true
				}
				return sqlTransactionListDatedToOnlyIncomeAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedToOnlyIncomeDescLimit, true
			}
			return sqlTransactionListDatedToOnlyIncomeDesc, true
		case sqlFilterAmountExpense:
			if oldestFirst {
				if withLimit {
					return sqlTransactionListDatedToOnlyExpenseAscLimit, true
				}
				return sqlTransactionListDatedToOnlyExpenseAsc, true
			}
			if withLimit {
				return sqlTransactionListDatedToOnlyExpenseDescLimit, true
			}
			return sqlTransactionListDatedToOnlyExpenseDesc, true
		}
	}
	return "", false
}

// staticListTransactionsSearchNoDateQuery returns fixed SQL for [ListTransactions] with search, no date bounds, optional LIMIT.
func staticListTransactionsSearchNoDateQuery(kind string, oldestFirst bool, limit int) (string, bool) {
	withLimit := limit > 0
	frag := sqlAmountKindFilter(kind)
	switch frag {
	case "":
		if oldestFirst {
			if withLimit {
				return sqlTransactionListSearchNoDateNoKindAscLimit, true
			}
			return sqlTransactionListSearchNoDateNoKindAsc, true
		}
		if withLimit {
			return sqlTransactionListSearchNoDateNoKindDescLimit, true
		}
		return sqlTransactionListSearchNoDateNoKindDesc, true
	case sqlFilterAmountIncome:
		if oldestFirst {
			if withLimit {
				return sqlTransactionListSearchNoDateIncomeAscLimit, true
			}
			return sqlTransactionListSearchNoDateIncomeAsc, true
		}
		if withLimit {
			return sqlTransactionListSearchNoDateIncomeDescLimit, true
		}
		return sqlTransactionListSearchNoDateIncomeDesc, true
	case sqlFilterAmountExpense:
		if oldestFirst {
			if withLimit {
				return sqlTransactionListSearchNoDateExpenseAscLimit, true
			}
			return sqlTransactionListSearchNoDateExpenseAsc, true
		}
		if withLimit {
			return sqlTransactionListSearchNoDateExpenseDescLimit, true
		}
		return sqlTransactionListSearchNoDateExpenseDesc, true
	}
	return "", false
}

// staticListTransactionsDatedBothSearchQuery returns fixed SQL when [ListTransactions] has both date bounds and a search term.
func staticListTransactionsDatedBothSearchQuery(kind string, oldestFirst bool, limit int) (string, bool) {
	withLimit := limit > 0
	frag := sqlAmountKindFilter(kind)
	switch frag {
	case "":
		if oldestFirst {
			if withLimit {
				return sqlTransactionListDatedBothSearchNoKindAscLimit, true
			}
			return sqlTransactionListDatedBothSearchNoKindAsc, true
		}
		if withLimit {
			return sqlTransactionListDatedBothSearchNoKindDescLimit, true
		}
		return sqlTransactionListDatedBothSearchNoKindDesc, true
	case sqlFilterAmountIncome:
		if oldestFirst {
			if withLimit {
				return sqlTransactionListDatedBothSearchIncomeAscLimit, true
			}
			return sqlTransactionListDatedBothSearchIncomeAsc, true
		}
		if withLimit {
			return sqlTransactionListDatedBothSearchIncomeDescLimit, true
		}
		return sqlTransactionListDatedBothSearchIncomeDesc, true
	case sqlFilterAmountExpense:
		if oldestFirst {
			if withLimit {
				return sqlTransactionListDatedBothSearchExpenseAscLimit, true
			}
			return sqlTransactionListDatedBothSearchExpenseAsc, true
		}
		if withLimit {
			return sqlTransactionListDatedBothSearchExpenseDescLimit, true
		}
		return sqlTransactionListDatedBothSearchExpenseDesc, true
	}
	return "", false
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
