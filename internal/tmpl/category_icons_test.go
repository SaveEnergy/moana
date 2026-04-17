package tmpl

import (
	"strings"
	"testing"

	"moana/internal/store"
)

func TestTemplateCategoryIconFormValue_legacyEmoji(t *testing.T) {
	t.Parallel()
	got := templateCategoryIconFormValue(store.Category{Icon: "💰"})
	if got != "wallet" {
		t.Fatalf("got %q want wallet", got)
	}
}

func TestCategoryNoneIconHTML(t *testing.T) {
	t.Parallel()
	h := categoryNoneIconHTML()
	s := string(h)
	if !strings.Contains(s, "<svg") || !strings.Contains(s, "moana-icon--sm") {
		t.Fatalf("got %q", s)
	}
}

func TestDashboardRecentIcon_namedCategoryUsesGlyph(t *testing.T) {
	t.Parallel()
	tx := store.Transaction{
		AmountCents:  -500,
		CategoryName: "Groceries",
		CategoryIcon: "utensils",
	}
	h := dashboardRecentIcon(tx, 0)
	if len(string(h)) < 30 || !strings.Contains(string(h), "<svg") {
		t.Fatalf("expected category svg, got len %d", len(string(h)))
	}
}

func TestDashboardRecentIcon_unsignedAmountTrendingUp(t *testing.T) {
	t.Parallel()
	tx := store.Transaction{AmountCents: 100, CategoryName: ""}
	h := dashboardRecentIcon(tx, 0)
	s := string(h)
	if !strings.Contains(s, "<svg") {
		t.Fatalf("got %q", s)
	}
}

func TestDashboardRecentIcon_expenseNoCategoryTrendingDown(t *testing.T) {
	t.Parallel()
	tx := store.Transaction{AmountCents: -100, CategoryName: ""}
	h := dashboardRecentIcon(tx, 0)
	s := string(h)
	if !strings.Contains(s, "<svg") {
		t.Fatalf("got %q", s)
	}
}
