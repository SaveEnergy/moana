package money

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseEURToCents parses a decimal euro amount (e.g. "1234.56", "1234") into integer cents.
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
	s = strings.ReplaceAll(s, ",", "")
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
	out := euros*100 + cents
	if neg {
		out = -out
	}
	return out, nil
}

// AbsCents returns the absolute value of an amount in cents.
func AbsCents(c int64) int64 {
	if c < 0 {
		return -c
	}
	return c
}
