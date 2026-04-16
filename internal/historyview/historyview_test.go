package historyview

import (
	"net/url"
	"testing"
	"time"
)

func TestBuildNav(t *testing.T) {
	t.Parallel()
	u, _ := url.Parse("/history?kind=expense&q=a")
	nav := BuildNav(u)
	if nav.LinkAll == "" || nav.SortOldest == "" {
		t.Fatal(nav)
	}
}

func TestBuildNav_preservesSearchAndDateFilters(t *testing.T) {
	t.Parallel()
	raw := "/history?kind=expense&q=coffee&from=2026-01-01&to=2026-01-31&sort=oldest"
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	nav := BuildNav(u)
	for _, link := range []string{nav.LinkAll, nav.LinkIncome, nav.LinkExpense, nav.SortNewest, nav.SortOldest} {
		parsed, err := url.Parse(link)
		if err != nil {
			t.Fatalf("parse %q: %v", link, err)
		}
		q := parsed.Query()
		if got := q.Get("q"); got != "coffee" {
			t.Fatalf("q: got %q want coffee (link=%s)", got, link)
		}
		if got := q.Get("from"); got != "2026-01-01" {
			t.Fatalf("from: got %q (link=%s)", got, link)
		}
		if got := q.Get("to"); got != "2026-01-31" {
			t.Fatalf("to: got %q (link=%s)", got, link)
		}
	}
	all, _ := url.Parse(nav.LinkAll)
	if all.Query().Get("kind") != "all" {
		t.Fatalf("LinkAll kind: %s", nav.LinkAll)
	}
	inc, _ := url.Parse(nav.LinkIncome)
	if inc.Query().Get("kind") != "income" {
		t.Fatalf("LinkIncome kind: %s", nav.LinkIncome)
	}
	exp, _ := url.Parse(nav.LinkExpense)
	if exp.Query().Get("kind") != "expense" {
		t.Fatalf("LinkExpense kind: %s", nav.LinkExpense)
	}
}

func TestBuildNav_sortNewestDropsSortParamKeepsSearch(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("/history?sort=oldest&q=rent")
	if err != nil {
		t.Fatal(err)
	}
	nav := BuildNav(u)
	parsed, err := url.Parse(nav.SortNewest)
	if err != nil {
		t.Fatal(err)
	}
	q := parsed.Query()
	if q.Get("sort") != "" {
		t.Fatalf("expected sort removed, got %q", q.Get("sort"))
	}
	if q.Get("q") != "rent" {
		t.Fatalf("expected q preserved, got %q", q.Get("q"))
	}
}

func TestFormatDayLabel_today(t *testing.T) {
	t.Parallel()
	loc := time.Local
	now := time.Now().In(loc)
	d0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	s := FormatDayLabel(d0, loc)
	if len(s) < 8 {
		t.Fatal(s)
	}
}
