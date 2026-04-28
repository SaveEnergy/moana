package store

import (
	"strings"
	"time"

	"moana/internal/timeutil"
)

// transactionListQueryPlan is the canonical SQL shape for [Store.ListTransactions].
//
// The previous implementation enumerated every date/search/kind/order/limit variant as a
// standalone const. This keeps the same fragment order but makes the filter matrix explicit:
// household -> date bounds -> amount kind -> search -> order -> limit.
type transactionListQueryPlan struct {
	fromUTC     bool
	toUTC       bool
	kindFilter  string
	search      bool
	oldestFirst bool
	limit       bool
}

func newTransactionListQueryPlan(fromUTC, toUTC *time.Time, kind, search string, oldestFirst bool, limit int) transactionListQueryPlan {
	return transactionListQueryPlan{
		fromUTC:     fromUTC != nil,
		toUTC:       toUTC != nil,
		kindFilter:  sqlAmountKindFilter(kind),
		search:      strings.TrimSpace(search) != "",
		oldestFirst: oldestFirst,
		limit:       limit > 0,
	}
}

func (p transactionListQueryPlan) SQL() string {
	var b strings.Builder
	b.Grow(512)
	b.WriteString(sqlTransactionListFromHousehold)
	if p.fromUTC {
		b.WriteString(sqlFilterOccurredAtFrom)
	}
	if p.toUTC {
		b.WriteString(sqlFilterOccurredAtTo)
	}
	if p.kindFilter != "" {
		b.WriteString(p.kindFilter)
	}
	if p.search {
		b.WriteString(sqlTransactionListSearchLike)
	}
	if p.oldestFirst {
		b.WriteString(sqlTransactionListOrderAsc)
	} else {
		b.WriteString(sqlTransactionListOrderDesc)
	}
	if p.limit {
		b.WriteString(sqlTransactionListLimitSuffix)
	}
	return b.String()
}

func transactionListQueryArgs(householdID int64, fromUTC, toUTC *time.Time, search string, limit int) []any {
	search = strings.TrimSpace(search)
	// At most: household, from, to, 2×search, limit.
	args := make([]any, 0, 6)
	args = append(args, householdID)
	if fromUTC != nil {
		args = append(args, timeutil.FormatSQLiteUTC(*fromUTC))
	}
	if toUTC != nil {
		args = append(args, timeutil.FormatSQLiteUTC(*toUTC))
	}
	if search != "" {
		term := "%" + escapeSQLLikePattern(search) + "%"
		args = append(args, term, term)
	}
	if limit > 0 {
		args = append(args, limit)
	}
	return args
}
