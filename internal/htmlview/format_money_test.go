package htmlview

import (
	"strings"
	"testing"

	"moana/internal/money"
)

func TestFormatEURAbs_negative(t *testing.T) {
	t.Parallel()
	got := FormatEURAbs(-12345)
	want := money.FormatEUR(12345)
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatCompactEUR_belowThresholdUsesFormatEUR(t *testing.T) {
	t.Parallel()
	got := FormatCompactEUR(99_999)
	want := money.FormatEUR(99_999)
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatCompactEUR_atThresholdAbbreviates(t *testing.T) {
	t.Parallel()
	if got := FormatCompactEUR(100_000); got != "€1.0k" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatCompactEUR_negativeLarge_preservesSign(t *testing.T) {
	t.Parallel()
	// -€1,500.00 → compact must not drop the minus
	got := FormatCompactEUR(-150_000)
	if !strings.HasPrefix(got, "-") {
		t.Fatalf("expected leading minus, got %q", got)
	}
	if !strings.Contains(got, "€") {
		t.Fatalf("got %q", got)
	}
}
