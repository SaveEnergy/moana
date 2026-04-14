package htmlview

import (
	"fmt"
	"html"
	"math"
	"time"

	"moana/internal/money"
)

// FormatEUR renders cents as EUR.
func FormatEUR(cents int64) string {
	return money.FormatEUR(cents)
}

// FormatRFC3339UTC formats t as RFC3339Nano in UTC for <time datetime> + client-side local display.
func FormatRFC3339UTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// Attr escapes text for safe use inside HTML double-quoted attributes.
func Attr(s string) string {
	return html.EscapeString(s)
}

// FormatEURAbs renders absolute cents as EUR.
func FormatEURAbs(cents int64) string {
	if cents < 0 {
		cents = -cents
	}
	return money.FormatEUR(cents)
}

// IsNegFloat reports whether x is a finite negative number.
func IsNegFloat(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0) && x < 0
}

// FormatPercentSigned formats a percentage with sign and one decimal.
func FormatPercentSigned(x float64) string {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return "—"
	}
	sign := ""
	if x >= 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%.1f%%", sign, x)
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
