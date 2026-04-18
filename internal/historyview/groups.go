package historyview

import (
	"sort"
	"strings"
	"time"

	"moana/internal/store"
	"moana/internal/timeutil"
)

// GroupByDay buckets transactions by local calendar day and orders day keys.
// A nil loc is treated as UTC ([time.Time.In] and [time.ParseInLocation] panic if loc is nil).
func GroupByDay(txs []store.Transaction, loc *time.Location, newestDayFirst bool) []DayGroup {
	if len(txs) == 0 {
		return nil
	}
	loc = timeutil.OrUTC(loc)
	// At most one bucket per row (distinct local days); hint reduces rehashing on large pages.
	byDay := make(map[string][]store.Transaction, len(txs))
	for _, tx := range txs {
		k := timeutil.LocalCalendarDateKey(tx.OccurredAt, loc)
		byDay[k] = append(byDay[k], tx)
	}
	keys := make([]string, 0, len(byDay))
	for k := range byDay {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if newestDayFirst {
			return keys[i] > keys[j]
		}
		return keys[i] < keys[j]
	})
	out := make([]DayGroup, 0, len(keys))
	for _, k := range keys {
		items := byDay[k]
		if len(items) == 0 {
			continue
		}
		// One midnight per bucket: keys come from local YYYY-MM-DD of these rows.
		t0 := items[0].OccurredAt.In(loc)
		midnight := time.Date(t0.Year(), t0.Month(), t0.Day(), 0, 0, 0, 0, loc)
		out = append(out, DayGroup{
			Label: FormatDayLabel(midnight, loc),
			Items: items,
		})
	}
	return out
}

// FormatDayLabel returns a short heading like "Today, JAN 2" or "Monday, Jan 2".
// A nil loc is treated as UTC.
func FormatDayLabel(day time.Time, loc *time.Location) string {
	loc = timeutil.OrUTC(loc)
	d := day.In(loc)
	now := time.Now().In(loc)
	d0 := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
	n0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	if d0.Equal(n0) {
		return "Today, " + strings.ToUpper(d.Format("Jan 2"))
	}
	y0 := n0.AddDate(0, 0, -1)
	if d0.Equal(y0) {
		return "Yesterday, " + strings.ToUpper(d.Format("Jan 2"))
	}
	return d.Format("Monday, Jan 2")
}
