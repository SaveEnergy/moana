package dashboard

import (
	"testing"

	"moana/internal/store"
)

func TestBuildOutflowFromExpenseRows_empty(t *testing.T) {
	t.Parallel()
	rows, donut, total := buildOutflowFromExpenseRows(nil, 0)
	if len(rows) != 0 || donut != "" || total != 0 {
		t.Fatalf("rows=%d donut=%q total=%d", len(rows), donut, total)
	}
}

func TestBuildOutflowFromExpenseRows_matchesPeriodExpenseDenominator(t *testing.T) {
	t.Parallel()
	// Single expense category; periodExpense from aggregate should match outflow denominator.
	in := []store.CategoryAmount{
		{Name: "Coffee", Icon: "cup", Color: "#112233", AmountCents: 5000},
	}
	outRows, donut, totalAbs := buildOutflowFromExpenseRows(in, -5000)
	if totalAbs != 5000 {
		t.Fatalf("OutflowTotalAbs %d want 5000", totalAbs)
	}
	if len(outRows) != 1 {
		t.Fatalf("len %d", len(outRows))
	}
	if outRows[0].Pct < 99.9 || outRows[0].Pct > 100.01 {
		t.Fatalf("Pct %f want ~100", outRows[0].Pct)
	}
	if donut == "" {
		t.Fatal("expected non-empty conic-gradient CSS")
	}
}
