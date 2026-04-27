package money

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ErrAmountTooLarge is returned when the parsed euro amount cannot fit in cents as int64.
var ErrAmountTooLarge = errors.New("amount too large")

// normalizeDecimalSeparators makes amounts parseable with the existing dot-based logic.
// It supports US (1,234.56) and common European (1.234,56 or 12,50) keypads where only a comma
// is available for the fractional part.
func normalizeDecimalSeparators(s string) string {
	lastDot := strings.LastIndex(s, ".")
	lastComma := strings.LastIndex(s, ",")
	if lastDot >= 0 && lastComma >= 0 {
		if lastComma > lastDot {
			// European: thousands are dots, decimal is the last comma.
			s = strings.ReplaceAll(s, ".", "")
			s = strings.ReplaceAll(s, ",", ".")
		} else {
			// US/UK: thousands are commas, dot is decimal.
			s = strings.ReplaceAll(s, ",", "")
		}
		return s
	}
	if lastComma >= 0 {
		switch strings.Count(s, ",") {
		case 1:
			frac := s[strings.Index(s, ",")+1:]
			if len(frac) <= 2 {
				// e.g. 12,50 or 0,5
				return strings.Replace(s, ",", ".", 1)
			}
			// e.g. 1,234 (thousands) or 12,123
			return strings.ReplaceAll(s, ",", "")
		default:
			// 1,234,567
			return strings.ReplaceAll(s, ",", "")
		}
	}
	// Dots only (or no separator): drop stray commas in US thousands form.
	return strings.ReplaceAll(s, ",", "")
}

// ParseEURToCents parses a decimal euro amount (e.g. "1234.56", "1.234,56", "12,50") into integer cents.
func ParseEURToCents(s string) (int64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "€", "")
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return 0, fmt.Errorf("amount is required")
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = strings.TrimPrefix(s, "-")
		s = strings.TrimSpace(s)
	}
	s = normalizeDecimalSeparators(s)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid amount")
	}
	var euros int64
	var err error
	if parts[0] == "" {
		euros = 0
	} else {
		euros, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount")
		}
	}
	var cents int64
	if len(parts) == 2 {
		frac := parts[1]
		if len(frac) > 2 {
			return 0, fmt.Errorf("use at most two decimal places")
		}
		for len(frac) < 2 {
			frac += "0"
		}
		c, err := strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount")
		}
		cents = c
	}
	// int64 multiply wraps on overflow; reject before euros*100+cents.
	maxEuros := int64(math.MaxInt64 / 100)
	maxRem := int64(math.MaxInt64 % 100)
	if euros > maxEuros || (euros == maxEuros && cents > maxRem) {
		return 0, ErrAmountTooLarge
	}
	out := euros*100 + cents
	if neg {
		out = -out
	}
	return out, nil
}

// absCentsUint64 returns the magnitude of c in cents. For [math.MinInt64], the true magnitude
// does not fit in int64; the result is uint64(math.MaxInt64)+1.
func absCentsUint64(c int64) uint64 {
	if c >= 0 {
		return uint64(c)
	}
	if c == math.MinInt64 {
		return uint64(math.MaxInt64) + 1
	}
	return uint64(-c)
}

// AbsCents returns the absolute value of an amount in cents. If the magnitude does not fit in
// int64 (only possible for [math.MinInt64]), it returns [math.MaxInt64] (saturation).
func AbsCents(c int64) int64 {
	u := absCentsUint64(c)
	if u > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(u)
}
