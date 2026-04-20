package dashboard

import (
	"math"

	"moana/internal/money"
	"moana/internal/store"
)

// pctChangeVsPrior is period-over-period % change: (current−previous) / |previous| × 100, with prior=0 handled as 0% or 100%.
// The difference uses float64 so current−previous does not wrap when it exceeds int64 range.
func pctChangeVsPrior(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	denom := float64(money.AbsCents(previous))
	numer := float64(current) - float64(previous)
	return numer / denom * 100
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
// If merging the tail would overflow int64, "Other" uses [math.MaxInt64] cents (saturation).
func MergeCategoryTopN(rows []store.CategoryAmount, limit int) []store.CategoryAmount {
	if limit < 1 || len(rows) <= limit {
		return rows
	}
	out := make([]store.CategoryAmount, limit)
	copy(out, rows[:limit-1])
	var rest int64
	for _, r := range rows[limit-1:] {
		var ok bool
		rest, ok = money.AddCents(rest, r.AmountCents)
		if !ok {
			rest = math.MaxInt64
			break
		}
	}
	out[limit-1] = store.CategoryAmount{Name: "Other", AmountCents: rest}
	return out
}
