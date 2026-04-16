package dashboard

import (
	"fmt"
	"strings"
)

// DonutConicGradient builds a CSS conic-gradient for expense slices (percentages sum to ~100).
func DonutConicGradient(pcts []float64, hexColors []string) string {
	n := len(pcts)
	if n == 0 {
		return ""
	}
	fallback := []string{"#306369", "#4a7d82", "#678a92", "#8aa3a8", "#b5c4c8"}
	var b strings.Builder
	// One segment is roughly " #rrggbb 12.345% 98.765%", comma-separated; avoid realloc while building.
	b.Grow(32 + n*80)
	b.WriteString("conic-gradient(from -90deg, ")
	cum := 0.0
	for i, p := range pcts {
		col := fallback[i%len(fallback)]
		if i < len(hexColors) && strings.TrimSpace(hexColors[i]) != "" {
			col = hexColors[i]
		}
		if i > 0 {
			b.WriteString(", ")
		}
		next := cum + p
		if next > 100.01 {
			next = 100
		}
		fmt.Fprintf(&b, "%s %.3f%% %.3f%%", col, cum, next)
		cum = next
	}
	b.WriteString(")")
	return b.String()
}
