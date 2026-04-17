package dashboard

import (
	"fmt"
	"math"
	"time"

	"moana/internal/money"
	"moana/internal/timeutil"
)

// HeatmapCell is one square in the year activity grid (padding cells use Empty=true).
type HeatmapCell struct {
	DateKey string
	Title   string
	Level   int
	Empty   bool
}

// BuildHeatmapCellsRolling365 builds a GitHub-style column-major grid (7 rows × N weeks) for the
// last 365 local calendar days ending on endDay (inclusive). The last cell is always endDay (today).
// Leading empty cells align the first day to the correct weekday row; no trailing padding.
// A nil loc is treated as UTC ([time.Date] panics if loc is nil).
func BuildHeatmapCellsRolling365(endDay time.Time, loc *time.Location, byDay map[string]int64) []HeatmapCell {
	loc = timeutil.OrUTC(loc)
	endDay = time.Date(endDay.Year(), endDay.Month(), endDay.Day(), 0, 0, 0, 0, loc)
	startDay := endDay.AddDate(0, 0, -364)

	maxC := int64(1)
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		k := d.Format("2006-01-02")
		if v := byDay[k]; v > maxC {
			maxC = v
		}
	}

	pad := int(startDay.Weekday())
	cells := make([]HeatmapCell, 0, pad+365)
	for i := 0; i < pad; i++ {
		cells = append(cells, HeatmapCell{Empty: true})
	}
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		cents := byDay[key]
		lvl := heatmapLevel(cents, maxC)
		title := fmt.Sprintf("%s — %s total activity", key, money.FormatEUR(cents))
		cells = append(cells, HeatmapCell{
			DateKey: key,
			Title:   title,
			Level:   lvl,
			Empty:   false,
		})
	}
	return cells
}

func heatmapLevel(cents, maxC int64) int {
	if cents <= 0 {
		return 0
	}
	if maxC <= 0 {
		return 1
	}
	x := float64(cents) / float64(maxC)
	lev := int(math.Ceil(x * 4))
	if lev < 1 {
		lev = 1
	}
	if lev > 4 {
		lev = 4
	}
	return lev
}
