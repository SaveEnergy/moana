package dashboard

import (
	"math"
	"strings"
	"testing"
	"time"

	"moana/internal/money"
	"moana/internal/store"
)

func TestParseStatsPeriod(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want statsPeriodConfig
	}{
		{"", statsPeriodConfig{Period: StatsPeriod30d, InclusiveDays: 30, PriorPhrase: "prior 30 days"}},
		{StatsPeriod30d, statsPeriodConfig{Period: StatsPeriod30d, InclusiveDays: 30, PriorPhrase: "prior 30 days"}},
		{StatsPeriod12m, statsPeriodConfig{Period: StatsPeriod12m, InclusiveDays: 365, PriorPhrase: "prior 12 months"}},
		{"garbage", statsPeriodConfig{Period: StatsPeriod30d, InclusiveDays: 30, PriorPhrase: "prior 30 days"}},
		// No trim: not exactly "12m", so treat like unknown (matches ?period= + accidental spaces in URL).
		{" 12m ", statsPeriodConfig{Period: StatsPeriod30d, InclusiveDays: 30, PriorPhrase: "prior 30 days"}},
	}
	for _, tc := range cases {
		got := parseStatsPeriod(tc.in)
		if got != tc.want {
			t.Fatalf("parseStatsPeriod(%q): %+v want %+v", tc.in, got, tc.want)
		}
	}
}

func TestPctChangeVsPrior_matchesExportedWrappers(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		cur, prev int64
	}{
		{150, 100},
		{0, 0},
		{10, 0},
		{0, 10},
		{100, -50},
		{0, -10},
	} {
		got := pctChangeVsPrior(tc.cur, tc.prev)
		if n := NetPctChange(tc.cur, tc.prev); got != n {
			t.Fatalf("cur=%d prev=%d: helper=%v NetPctChange=%v", tc.cur, tc.prev, got, n)
		}
		if p := PctChangePositive(tc.cur, tc.prev); got != p {
			t.Fatalf("cur=%d prev=%d: helper=%v PctChangePositive=%v", tc.cur, tc.prev, got, p)
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

func TestNetPctChange_negativePreviousUsesAbsInDenominator(t *testing.T) {
	t.Parallel()
	// current 100, previous -50 -> (100-(-50)) / Abs(-50) * 100 = 300
	v := NetPctChange(100, -50)
	if v < 299.9 || v > 300.1 {
		t.Fatalf("got %v want ~300", v)
	}
}

func TestNetPctChange_extremePriorMinInt64_noInt64SubtractionWrap(t *testing.T) {
	t.Parallel()
	// int64(100 - MinInt64) wraps to a negative value; the % must use float difference (~+100%).
	v := NetPctChange(100, math.MinInt64)
	if v < 99 || v > 101 {
		t.Fatalf("got %v want ~100 (wrapped int64 diff would yield ~-100)", v)
	}
}

func TestNetPctChange_maxVsMinInt64_nearTwoHundredPercent(t *testing.T) {
	t.Parallel()
	v := NetPctChange(math.MaxInt64, math.MinInt64)
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Fatalf("got %v", v)
	}
	// (Max−Min) / |Min| ≈ 2 → ~200% (int64 subtraction wraps to −1 and would be nonsense).
	if v < 199 || v > 201 {
		t.Fatalf("got %v want ~200", v)
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

func TestPctChangePositive_negativePreviousUsesAbsDenominator(t *testing.T) {
	t.Parallel()
	// Prior total was negative (bad aggregate); denominator uses magnitude like NetPctChange.
	// (0-(-10))/10*100 = 100
	v := PctChangePositive(0, -10)
	if v < 99.9 || v > 100.1 {
		t.Fatalf("got %v want ~100", v)
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

// TestMergeCategoryTopN_limitOneRollsAllIntoOther documents limit=1: no named top rows, single "Other" bucket.
func TestMergeCategoryTopN_sumConservedWhenMerging(t *testing.T) {
	t.Parallel()
	rows := []store.CategoryAmount{
		{Name: "a", AmountCents: 100},
		{Name: "b", AmountCents: 200},
		{Name: "c", AmountCents: 300},
		{Name: "d", AmountCents: 400},
	}
	var total int64
	for _, r := range rows {
		total += r.AmountCents
	}
	out := MergeCategoryTopN(rows, 3)
	var got int64
	for _, r := range out {
		got += r.AmountCents
	}
	if got != total {
		t.Fatalf("sum %d want %d (out=%+v)", got, total, out)
	}
	if len(out) != 3 || out[2].Name != "Other" || out[2].AmountCents != 700 {
		t.Fatalf("got %+v", out)
	}
}

func TestMergeCategoryTopN_limitOneRollsAllIntoOther(t *testing.T) {
	t.Parallel()
	rows := []store.CategoryAmount{
		{Name: "a", AmountCents: 100},
		{Name: "b", AmountCents: 200},
	}
	out := MergeCategoryTopN(rows, 1)
	if len(out) != 1 || out[0].Name != "Other" || out[0].AmountCents != 300 {
		t.Fatalf("got %+v", out)
	}
}

func TestMergeCategoryTopN_otherSaturatesWhenSumOverflows(t *testing.T) {
	t.Parallel()
	rows := []store.CategoryAmount{
		{Name: "top", AmountCents: 1},
		{Name: "b", AmountCents: math.MaxInt64},
		{Name: "c", AmountCents: math.MaxInt64},
	}
	out := MergeCategoryTopN(rows, 2)
	if len(out) != 2 || out[1].Name != "Other" {
		t.Fatalf("got %+v", out)
	}
	if out[1].AmountCents != math.MaxInt64 {
		t.Fatalf("Other bucket want MaxInt64 when merge sum overflows int64, got %d", out[1].AmountCents)
	}
	if _, ok := money.AddCents(rows[1].AmountCents, rows[2].AmountCents); ok {
		t.Fatal("test setup: expected raw sum to overflow")
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

func TestBuildHeatmapCellsRolling365_nilLocationMatchesUTC(t *testing.T) {
	t.Parallel()
	end := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)
	byDay := map[string]int64{"2026-04-13": 100}
	got := BuildHeatmapCellsRolling365(end, nil, byDay)
	want := BuildHeatmapCellsRolling365(end, time.UTC, byDay)
	if len(got) != len(want) {
		t.Fatalf("len %d vs %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("cell %d: %+v vs %+v", i, got[i], want[i])
		}
	}
}

func TestDonutConicGradient(t *testing.T) {
	t.Parallel()
	s := DonutConicGradient([]float64{40, 60}, []string{"#111111", "#222222"})
	if s == "" || len(s) < 20 {
		t.Fatal(s)
	}
	if !strings.HasPrefix(s, "conic-gradient(from -90deg, ") {
		t.Fatalf("expected CSS conic-gradient prefix, got %q", s)
	}
	if !strings.HasSuffix(s, ")") {
		t.Fatalf("expected closed gradient, got %q", s)
	}
}

func TestDonutConicGradient_empty(t *testing.T) {
	t.Parallel()
	if DonutConicGradient(nil, nil) != "" {
		t.Fatal("nil slices")
	}
	if DonutConicGradient([]float64{}, []string{}) != "" {
		t.Fatal("empty slices")
	}
}

func TestDonutConicGradient_fallbackColorWhenHexBlank(t *testing.T) {
	t.Parallel()
	got := DonutConicGradient([]float64{50, 50}, []string{"#111111", "  "})
	if !strings.Contains(got, "#111111") || !strings.Contains(got, "#4a7d82") {
		t.Fatalf("expected custom first + palette fallback second, got %q", got)
	}
}

func TestDonutConicGradient_clampsCumulativePast100(t *testing.T) {
	t.Parallel()
	got := DonutConicGradient([]float64{60, 50}, []string{"#aaaaaa", "#bbbbbb"})
	if !strings.Contains(got, "100.000%") {
		t.Fatalf("expected cap at 100%%, got %q", got)
	}
}
