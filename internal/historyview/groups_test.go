package historyview

import (
	"testing"
	"time"

	"moana/internal/store"
)

func TestGroupByDay_empty(t *testing.T) {
	t.Parallel()
	if g := GroupByDay(nil, time.UTC, true); g != nil {
		t.Fatalf("got %#v", g)
	}
	if g := GroupByDay([]store.Transaction{}, time.UTC, true); g != nil {
		t.Fatalf("got %#v", g)
	}
}

func TestGroupByDay_bucketsAndOrder(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	day1a := time.Date(2026, 1, 2, 10, 0, 0, 0, loc)
	day1b := time.Date(2026, 1, 2, 15, 30, 0, 0, loc)
	day2 := time.Date(2026, 1, 3, 9, 0, 0, 0, loc)
	txs := []store.Transaction{
		{OccurredAt: day1a, Description: "morning"},
		{OccurredAt: day1b, Description: "afternoon"},
		{OccurredAt: day2, Description: "next"},
	}
	g := GroupByDay(txs, loc, true)
	if len(g) != 2 {
		t.Fatalf("groups %d", len(g))
	}
	// Newest day first: 2026-01-03 then 2026-01-02
	if len(g[0].Items) != 1 || g[0].Items[0].Description != "next" {
		t.Fatalf("first group: %+v", g[0])
	}
	if len(g[1].Items) != 2 {
		t.Fatalf("second group len %d", len(g[1].Items))
	}
	gOld := GroupByDay(txs, loc, false)
	if len(gOld) != 2 || len(gOld[0].Items) != 2 {
		t.Fatalf("oldest-first first group: %+v", gOld[0])
	}
}
