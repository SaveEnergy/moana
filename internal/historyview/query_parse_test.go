package historyview

import (
	"net/url"
	"testing"
)

func TestParseHistoryURL_defaults(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("/history")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != "all" || p.filterKind != "" || p.search != "" || p.sortLabel != "newest" || p.oldestFirst || p.filterActive {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_kindAndSort(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("/history?kind=income&sort=oldest&q=rent")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != "income" || p.filterKind != "income" || p.search != "rent" || p.sortLabel != "oldest" || !p.oldestFirst {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_dateFilterActive(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("/history?from=2026-01-01&to=2026-01-31")
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if !p.filterActive || p.from != "2026-01-01" || p.to != "2026-01-31" {
		t.Fatalf("%+v", p)
	}
}

func TestParseHistoryURL_trimsDateFields(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("/history?from=%20%202026-01-01%20%20&to=2026-01-31")
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
	if historyReturnOrDefault("") != "/history" {
		t.Fatal()
	}
	if historyReturnOrDefault("/history?q=a") != "/history?q=a" {
		t.Fatal()
	}
}
