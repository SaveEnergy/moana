package money

import "testing"

func TestFormatEUR(t *testing.T) {
	t.Parallel()
	tests := []struct {
		cents int64
		want  string
	}{
		{0, "€0.00"},
		{1, "€0.01"},
		{99, "€0.99"},
		{12345, "€123.45"},
		{100000, "€1,000.00"},
		{123456789, "€1,234,567.89"},
		{99999999999999900, "€999,999,999,999,999.00"},
		{-500, "-€5.00"},
	}
	for _, tc := range tests {
		if got := FormatEUR(tc.cents); got != tc.want {
			t.Errorf("FormatEUR(%d) = %q, want %q", tc.cents, got, tc.want)
		}
	}
}

func TestFormatDecimalEURAbs(t *testing.T) {
	t.Parallel()
	if got := FormatDecimalEURAbs(12345); got != "123.45" {
		t.Errorf("got %q", got)
	}
	if got := FormatDecimalEURAbs(-9900); got != "99.00" {
		t.Errorf("got %q", got)
	}
}
