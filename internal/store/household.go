package store

import "time"

// Household is a shared account group.
type Household struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}
