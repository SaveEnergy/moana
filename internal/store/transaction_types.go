package store

import (
	"database/sql"
	"time"
)

// Transaction is a single income (or adjustment) entry.
type Transaction struct {
	ID           int64
	UserID       int64
	AmountCents  int64
	OccurredAt   time.Time
	Description  string
	CategoryID   sql.NullInt64
	CategoryName string
	CategoryIcon string // from categories.icon when joined
	CreatedAt    time.Time
}

// TransactionFilter restricts listed transactions.
type TransactionFilter struct {
	FromUTC     *time.Time
	ToUTC       *time.Time
	Kind        string // "", "all", "income", "expense"
	Search      string // matches description or category name (substring)
	OldestFirst bool   // default false = newest first
	Limit       int    // if > 0, cap rows (newest/oldest order preserved)
}
