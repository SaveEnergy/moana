package money

import (
	"fmt"
	"strings"
)

// FormatEUR formats cents as English EUR (e.g. €1,234.56). Negative amounts show a leading minus.
func FormatEUR(cents int64) string {
	neg := cents < 0
	abs := absCentsUint64(cents)
	whole := abs / 100
	frac := abs % 100
	intStr := formatThousandsUint64(whole)
	s := "€" + intStr + fmt.Sprintf(".%02d", frac)
	if neg {
		return "-" + s
	}
	return s
}

// FormatDecimalEURAbs formats absolute cents as a plain decimal (e.g. "1234.56") for HTML inputs.
func FormatDecimalEURAbs(cents int64) string {
	abs := absCentsUint64(cents)
	return fmt.Sprintf("%d.%02d", abs/100, abs%100)
}

func formatThousandsUint64(n uint64) string {
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
