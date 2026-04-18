package dashboard

import (
	"testing"
	"time"
)

func TestHeatmapLevel_zeroOrNegativeCents(t *testing.T) {
	t.Parallel()
	if heatmapLevel(0, 100) != 0 {
		t.Fatalf("got %d", heatmapLevel(0, 100))
	}
	if heatmapLevel(-10, 100) != 0 {
		t.Fatalf("got %d", heatmapLevel(-10, 100))
	}
}

func TestHeatmapLevel_nonPositiveMax(t *testing.T) {
	t.Parallel()
	// maxC<=0 falls back to level 1 when there is activity (see heatmap.go)
	if heatmapLevel(50, 0) != 1 {
		t.Fatalf("got %d", heatmapLevel(50, 0))
	}
	if heatmapLevel(50, -1) != 1 {
		t.Fatalf("got %d", heatmapLevel(50, -1))
	}
}

func TestHeatmapLevel_bucketsAgainstMax(t *testing.T) {
	t.Parallel()
	const maxC int64 = 100
	if got := heatmapLevel(25, maxC); got != 1 {
		t.Fatalf("25/100 -> %d want 1", got)
	}
	if got := heatmapLevel(50, maxC); got != 2 {
		t.Fatalf("50/100 -> %d want 2", got)
	}
	if got := heatmapLevel(100, maxC); got != 4 {
		t.Fatalf("100/100 -> %d want 4", got)
	}
}

func TestHeatmapLevel_clampedToFour(t *testing.T) {
	t.Parallel()
	if got := heatmapLevel(200, 100); got != 4 {
		t.Fatalf("cents>maxC -> %d want 4", got)
	}
}

func TestBuildHeatmapCellsRolling365_lenPaddingAndLastDayUTC(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	end := time.Date(2026, 1, 10, 15, 30, 0, 0, loc)
	start := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -364)
	pad := int(start.Weekday())
	cells := BuildHeatmapCellsRolling365(end, loc, nil)
	if len(cells) != pad+365 {
		t.Fatalf("len=%d want pad+365=%d", len(cells), pad+365)
	}
	last := cells[len(cells)-1]
	if last.Empty {
		t.Fatal("last cell must not be padding")
	}
	if last.DateKey != "2026-01-10" {
		t.Fatalf("last DateKey=%q want 2026-01-10", last.DateKey)
	}
}

func TestBuildHeatmapCellsRolling365_nonzeroDayUsesByDay(t *testing.T) {
	t.Parallel()
	loc := time.UTC
	end := time.Date(2026, 3, 1, 0, 0, 0, 0, loc)
	byDay := map[string]int64{"2026-03-01": 10_000}
	cells := BuildHeatmapCellsRolling365(end, loc, byDay)
	last := cells[len(cells)-1]
	if last.Empty || last.Level == 0 {
		t.Fatalf("expected activity cell, got %+v", last)
	}
	if last.Title == "" {
		t.Fatal("title must include formatted money")
	}
}
