package store

import "strings"

// SQL AND fragments for filtering ledger rows by amount sign. Used by listing and aggregates
// so "income" / "expense" semantics stay aligned across queries.
const (
	sqlFilterAmountIncome  = ` AND t.amount_cents > 0`
	sqlFilterAmountExpense = ` AND t.amount_cents < 0`
)

// sqlAmountKindFilter returns sqlFilterAmountIncome, sqlFilterAmountExpense, or "" for unknown kind.
// Leading and trailing ASCII space are ignored so callers match [strings.TrimSpace] behavior used by listing.
func sqlAmountKindFilter(kind string) string {
	kind = strings.TrimSpace(kind)
	switch kind {
	case "income":
		return sqlFilterAmountIncome
	case "expense":
		return sqlFilterAmountExpense
	default:
		return ""
	}
}
