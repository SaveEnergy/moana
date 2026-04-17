package txform

import (
	"testing"
	"time"

	"moana/internal/money"
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

func TestParse_amountTooLarge(t *testing.T) {
	t.Parallel()
	_, errMsg := Parse("922337203685477580.00", "2026-01-01", "", "", "income", time.UTC)
	if errMsg != money.ErrAmountTooLarge.Error() {
		t.Fatalf("got errMsg %q want overflow from money.ParseEURToCents", errMsg)
	}
}

func TestParse_incomePositiveCents(t *testing.T) {
	t.Parallel()
	p, errMsg := Parse("42.50", "2026-06-15", "bonus", "", "income", time.UTC)
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	if p.AmountCents != 4250 || p.Description != "bonus" {
		t.Fatalf("%+v", p)
	}
}

func TestParse_invalidDate(t *testing.T) {
	t.Parallel()
	_, errMsg := Parse("10", "not-a-date", "x", "", "expense", time.UTC)
	if errMsg == "" {
		t.Fatal("expected error")
	}
}

func TestParse_emptyDate(t *testing.T) {
	t.Parallel()
	_, errMsg := Parse("10.00", "", "x", "", "income", time.UTC)
	if errMsg != "Date is required." {
		t.Fatalf("got %q", errMsg)
	}
}

func TestParse_zeroAmount(t *testing.T) {
	t.Parallel()
	_, errMsg := Parse("0", "2026-06-15", "x", "", "income", time.UTC)
	if errMsg != "Amount must be greater than zero." {
		t.Fatalf("got %q", errMsg)
	}
}

func TestParse_invalidCategoryID(t *testing.T) {
	t.Parallel()
	_, errMsg := Parse("10.00", "2026-06-15", "x", "not-an-int", "expense", time.UTC)
	if errMsg != "Invalid category." {
		t.Fatalf("got %q", errMsg)
	}
}

func TestParse_categoryIdZero(t *testing.T) {
	t.Parallel()
	p, errMsg := Parse("10.00", "2026-06-15", "x", "0", "expense", time.UTC)
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	if p.CategoryID == nil || *p.CategoryID != 0 {
		t.Fatalf("got %+v", p)
	}
}

func TestParse_descriptionTrimmed(t *testing.T) {
	t.Parallel()
	p, errMsg := Parse("10.00", "2026-06-15", "  coffee  ", "", "expense", time.UTC)
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	if p.Description != "coffee" {
		t.Fatalf("got %q", p.Description)
	}
}

func TestParse_unknownKindDefaultsToIncome(t *testing.T) {
	t.Parallel()
	p, errMsg := Parse("5.00", "2026-06-15", "x", "", "weird", time.UTC)
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	if p.AmountCents != 500 {
		t.Fatalf("want positive income cents, got %+v", p)
	}
}

func TestParse_nilLocationUsesUTC(t *testing.T) {
	t.Parallel()
	p, errMsg := Parse("1.00", "2026-03-20", "", "", "income", nil)
	if errMsg != "" {
		t.Fatal(errMsg)
	}
	want := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
	if !p.OccurredUTC.Equal(want) {
		t.Fatalf("OccurredUTC=%v want %v", p.OccurredUTC, want)
	}
}
