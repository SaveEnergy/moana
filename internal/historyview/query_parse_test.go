package historyview

import (
	"net/url"
	"testing"
)

func TestParseHistoryURLValues_matchesParseHistoryURL(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?kind=expense&from=2026-01-01&to=2026-01-31&q=rent")
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if parseHistoryURLValues(q) != ParseHistoryURL(u) {
		t.Fatalf("parseHistoryURLValues vs ParseHistoryURL")
	}
}

func TestParseHistoryURL_defaults(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath)
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != KindAll || p.filterKind != "" || p.search != "" || p.sortLabel != SortLabelNewest || p.oldestFirst || p.filterActive {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_kindAndSort(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?kind=income&sort=oldest&q=rent")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != KindIncome || p.filterKind != KindIncome || p.search != "rent" || p.sortLabel != SortOldestValue || !p.oldestFirst {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_unknownKindFallsBackToAll(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?kind=nope&sort=oldest")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != KindAll || p.filterKind != "" || !p.oldestFirst || p.sortLabel != SortOldestValue {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_unknownSortFallsBackToNewest(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?kind=expense&sort=newestish")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != KindExpense || p.filterKind != KindExpense || p.oldestFirst || p.sortLabel != SortLabelNewest {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_dateFilterActive(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?from=2026-01-01&to=2026-01-31")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if !p.filterActive || p.from != "2026-01-01" || p.to != "2026-01-31" {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_allQueryParamsTogether(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?kind=expense&sort=oldest&q=coffee&from=2026-01-01&to=2026-01-31")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != KindExpense || p.filterKind != KindExpense || p.search != "coffee" ||
		p.sortLabel != SortOldestValue || !p.oldestFirst || !p.filterActive ||
		p.from != "2026-01-01" || p.to != "2026-01-31" {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_trimsDateFields(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?from=%20%202026-01-01%20%20&to=2026-01-31")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.from != "2026-01-01" || p.to != "2026-01-31" {
		t.Fatalf("%+v", p)
	}
}

func TestHistoryReturnOrDefault(t *testing.T) {
	t.Parallel()
	if historyReturnOrDefault("") != RoutePath {
		t.Fatal()
	}
	if historyReturnOrDefault(RoutePath+"?q=a") != RoutePath+"?q=a" {
		t.Fatal()
	}
}
