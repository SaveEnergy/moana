package money

import (
	"fmt"
	"strings"
)

// FormatEUR formats cents as English EUR (e.g. €1,234.56). Negative amounts show a leading minus.
func FormatEUR(cents int64) string {
	neg := cents < 0
	if neg {
		cents = -cents
	}
	whole := cents / 100
	frac := cents % 100
	intStr := formatThousands(whole)
	s := "€" + intStr + fmt.Sprintf(".%02d", frac)
	if neg {
		return "-" + s
	}
	return s
}

// FormatDecimalEURAbs formats absolute cents as a plain decimal (e.g. "1234.56") for HTML inputs.
func FormatDecimalEURAbs(cents int64) string {
	if cents < 0 {
		cents = -cents
	}
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

func formatThousands(n int64) string {
	if n < 0 {
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + (len(s)-1)/3)
	lead := len(s) % 3
	if lead == 0 {
		lead = 3
	}
	b.WriteString(s[:lead])
	for i := lead; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}
	return b.String()
}
