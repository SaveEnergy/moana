package store

// Category is a household-owned label for transactions (shared by all members).
type Category struct {
	ID          int64
	HouseholdID int64
	Name        string
	Icon        string // Lucide icon id from app palette; empty = heuristic in UI
	Color       string // CSS hex accent, e.g. #ea580c; empty = computed accent in UI
}
