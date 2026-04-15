package htmlview

import (
	"fmt"
	"html"
	"math"
	"time"

	"moana/internal/timeutil"
)

// FormatRFC3339UTC formats t as RFC3339Nano in UTC for <time datetime> + client-side local display.
func FormatRFC3339UTC(t time.Time) string {
	return timeutil.FormatSQLiteUTC(t)
}

// Attr escapes text for safe use inside HTML double-quoted attributes.
func Attr(s string) string {
	return html.EscapeString(s)
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
