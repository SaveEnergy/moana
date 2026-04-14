package txform

import (
	"strconv"
	"strings"
	"time"

	"moana/internal/money"
)

// Parsed is the normalized result of a transaction form POST.
type Parsed struct {
	AmountCents int64
	OccurredUTC time.Time
	Description string
	CategoryID  *int64
}

// Parse converts form strings into cents, UTC time, and optional category.
// kind is typically "income" or "expense" from the form; loc is the user's display timezone for the date field.
// On validation failure, errMsg is a short English message suitable for the UI.
func Parse(amountStr, dateStr, desc, catStr, kind string, loc *time.Location) (Parsed, string) {
	cents, err := money.ParseEURToCents(amountStr)
	if err != nil {
		return Parsed{}, err.Error()
	}
	cents = money.AbsCents(cents)
	if cents == 0 {
		return Parsed{}, "Amount must be greater than zero."
	}
	if kind == "expense" {
		cents = -cents
	}
	if strings.TrimSpace(dateStr) == "" {
		return Parsed{}, "Date is required."
	}
	dayStart, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return Parsed{}, "Invalid date."
	}
	var catID *int64
	if catStr != "" {
		id, err := strconv.ParseInt(catStr, 10, 64)
		if err != nil {
			return Parsed{}, "Invalid category."
		}
		catID = &id
	}
	return Parsed{
		AmountCents: cents,
		OccurredUTC: dayStart.UTC(),
		Description: desc,
		CategoryID:  catID,
	}, ""
}
