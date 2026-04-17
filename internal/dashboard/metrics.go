package dashboard

import (
	"moana/internal/money"
	"moana/internal/store"
)

// NetPctChange is period-over-period % change for signed net (current vs previous period of same length).
func NetPctChange(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return float64(current-previous) / float64(money.AbsCents(previous)) * 100
}

// PctChangePositive is period-over-period % change for non-negative amounts (income totals or expense absolutes).
// The denominator uses [money.AbsCents] so a negative prior total (unexpected input) does not flip the sign of the ratio.
func PctChangePositive(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return float64(current-previous) / float64(money.AbsCents(previous)) * 100
}

// MergeCategoryTopN keeps the top (limit-1) categories and merges the rest into "Other".
// If limit is less than 1, rows are returned unchanged (defensive; production uses a fixed positive limit).
func MergeCategoryTopN(rows []store.CategoryAmount, limit int) []store.CategoryAmount {
	if limit < 1 || len(rows) <= limit {
		return rows
	}
	out := make([]store.CategoryAmount, limit)
	copy(out, rows[:limit-1])
	var rest int64
	for _, r := range rows[limit-1:] {
		rest += r.AmountCents
	}
	out[limit-1] = store.CategoryAmount{Name: "Other", AmountCents: rest}
	return out
}
