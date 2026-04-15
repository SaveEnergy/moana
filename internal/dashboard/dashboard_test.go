package dashboard

import (
	"testing"
	"time"

	"moana/internal/store"
)

func TestParseStatsPeriod(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want statsPeriodConfig
	}{
		{"", statsPeriodConfig{Period: "30d", InclusiveDays: 30, PriorPhrase: "prior 30 days"}},
		{"30d", statsPeriodConfig{Period: "30d", InclusiveDays: 30, PriorPhrase: "prior 30 days"}},
		{"12m", statsPeriodConfig{Period: "12m", InclusiveDays: 365, PriorPhrase: "prior 12 months"}},
		{"garbage", statsPeriodConfig{Period: "30d", InclusiveDays: 30, PriorPhrase: "prior 30 days"}},
	}
	for _, tc := range cases {
		got := parseStatsPeriod(tc.in)
		if got != tc.want {
			t.Fatalf("parseStatsPeriod(%q): %+v want %+v", tc.in, got, tc.want)
		}
	}
}

func TestNetPctChange(t *testing.T) {
	t.Parallel()
	if v := NetPctChange(150, 100); v < 49.9 || v > 50.1 {
		t.Fatalf("got %v", v)
	}
	if v := NetPctChange(-50, -100); v < 49.9 || v > 50.1 {
		t.Fatalf("got %v", v)
	}
	if v := NetPctChange(0, 0); v != 0 {
		t.Fatalf("got %v", v)
	}
}

func TestPctChangePositive(t *testing.T) {
	t.Parallel()
	if v := PctChangePositive(150, 100); v < 49.9 || v > 50.1 {
		t.Fatalf("got %v", v)
	}
	if v := PctChangePositive(0, 0); v != 0 {
		t.Fatalf("both zero: got %v", v)
	}
	if v := PctChangePositive(10, 0); v != 100 {
		t.Fatalf("prior zero: got %v want 100", v)
	}
	if v := PctChangePositive(0, 10); v != -100 {
		t.Fatalf("current zero: got %v want -100", v)
	}
}

func TestMergeCategoryTopN(t *testing.T) {
	t.Parallel()
	rows := []store.CategoryAmount{
		{Name: "a", AmountCents: 100},
		{Name: "b", AmountCents: 200},
		{Name: "c", AmountCents: 300},
	}
	out := MergeCategoryTopN(rows, 2)
	if len(out) != 2 || out[0].Name != "a" || out[1].Name != "Other" || out[1].AmountCents != 500 {
		t.Fatalf("got %+v", out)
	}
}

func TestMergeCategoryTopN_invalidLimitReturnsRows(t *testing.T) {
	t.Parallel()
	rows := []store.CategoryAmount{
		{Name: "a", AmountCents: 100},
		{Name: "b", AmountCents: 200},
	}
	out := MergeCategoryTopN(rows, 0)
	if len(out) != 2 || out[0].Name != "a" {
		t.Fatalf("got %+v", out)
	}
}

func TestBuildHeatmapCellsRolling365_todayIsLastCell(t *testing.T) {
	t.Parallel()
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	end := time.Date(2026, 4, 13, 0, 0, 0, 0, loc)
	byDay := map[string]int64{
		"2025-04-13": 100,
		"2026-04-13": 200,
	}
	cells := BuildHeatmapCellsRolling365(end, loc, byDay)
	var lastNonEmpty string
	for i := len(cells) - 1; i >= 0; i-- {
		if !cells[i].Empty {
			lastNonEmpty = cells[i].DateKey
			break
		}
	}
	if lastNonEmpty != "2026-04-13" {
		t.Fatalf("last data cell want 2026-04-13, got %q", lastNonEmpty)
	}
	if len(cells) != 365+int(end.AddDate(0, 0, -364).Weekday()) {
		t.Fatalf("unexpected cell count: %d", len(cells))
	}
}

func TestDonutConicGradient(t *testing.T) {
	t.Parallel()
	s := DonutConicGradient([]float64{40, 60}, []string{"#111111", "#222222"})
	if s == "" || len(s) < 20 {
		t.Fatal(s)
	}
}
