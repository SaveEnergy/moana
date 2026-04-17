package dashboard

import "testing"

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
