package dashboard

import (
	"moana/internal/money"
	"moana/internal/store"
)

// pctChangeVsPrior is period-over-period % change: (current−previous) / |previous| × 100, with prior=0 handled as 0% or 100%.
func pctChangeVsPrior(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return float64(current-previous) / float64(money.AbsCents(previous)) * 100
}

// NetPctChange is period-over-period % change for signed net (current vs previous period of same length).
func NetPctChange(current, previous int64) float64 {
	return pctChangeVsPrior(current, previous)
}

// PctChangePositive is period-over-period % change for income totals or expense absolutes (same formula as [NetPctChange]; name reflects call-site intent).
func PctChangePositive(current, previous int64) float64 {
	return pctChangeVsPrior(current, previous)
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
