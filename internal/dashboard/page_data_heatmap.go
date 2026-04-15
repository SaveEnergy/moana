package dashboard

import (
	"context"
	"fmt"
	"time"

	"moana/internal/store"
)

func buildHeatmapSection(ctx context.Context, st *store.Store, householdID int64, loc *time.Location, now time.Time) (string, []HeatmapCell, int, error) {
	todayLocal := time.Date(now.In(loc).Year(), now.In(loc).Month(), now.In(loc).Day(), 0, 0, 0, 0, loc)
	yearAgo := todayLocal.AddDate(0, 0, -364)
	startUTC := yearAgo.UTC()
	endOfToday := time.Date(todayLocal.Year(), todayLocal.Month(), todayLocal.Day(), 23, 59, 59, 999999999, todayLocal.Location())
	endUTC := endOfToday.UTC()
	heatmapRangeLabel := fmt.Sprintf("%s – %s", yearAgo.Format("Jan 2, 2006"), todayLocal.Format("Jan 2, 2006"))

	byDay, err := st.DailyAbsMovementByLocalDate(ctx, householdID, startUTC, endUTC, loc)
	if err != nil {
		return "", nil, 0, err
	}
	heatmapCells := BuildHeatmapCellsRolling365(todayLocal, loc, byDay)
	heatmapCols := 1
	if n := len(heatmapCells); n > 0 {
		heatmapCols = (n + 6) / 7
	}
	return heatmapRangeLabel, heatmapCells, heatmapCols, nil
}
