package txform

import (
	"testing"
	"time"
)

func TestParse_ok(t *testing.T) {
	t.Parallel()
	p, errMsg := Parse("10.00", "2026-06-15", "x", "", "expense", time.UTC)
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	if p.AmountCents != -1000 || p.Description != "x" || p.CategoryID != nil {
		t.Fatalf("%+v", p)
	}
}

func TestParse_invalidAmount(t *testing.T) {
	t.Parallel()
	_, errMsg := Parse("abc", "2026-01-01", "", "", "income", time.UTC)
	if errMsg == "" {
		t.Fatal("expected error")
	}
}
