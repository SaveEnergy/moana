package txform

import "time"

// Parsed is the normalized result of a transaction form POST.
type Parsed struct {
	AmountCents int64
	OccurredUTC time.Time
	Description string
	CategoryID  *int64
}
