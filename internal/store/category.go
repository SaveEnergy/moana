package store

// Category is a user-owned label for transactions.
type Category struct {
	ID     int64
	UserID int64
	Name   string
	Icon   string // Lucide icon id from app palette; empty = heuristic in UI
	Color  string // CSS hex accent, e.g. #ea580c; empty = computed accent in UI
}
