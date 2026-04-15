package htmlview

import (
	"fmt"

	"moana/internal/money"
)

// FormatEUR renders cents as EUR.
func FormatEUR(cents int64) string {
	return money.FormatEUR(cents)
}

// FormatEURAbs renders absolute cents as EUR.
func FormatEURAbs(cents int64) string {
	if cents < 0 {
		cents = -cents
	}
	return money.FormatEUR(cents)
}

// FormatCompactEUR abbreviates large EUR amounts (e.g. €12.5k).
func FormatCompactEUR(cents int64) string {
	x := cents
	if x < 0 {
		x = -x
	}
	if x < 100_000 {
		return money.FormatEUR(cents)
	}
	v := float64(x) / 100.0 / 1000.0
	return fmt.Sprintf("€%.1fk", v)
}
