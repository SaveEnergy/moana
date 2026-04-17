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
	byDay := make(map[string][]store.Transaction)
	for _, tx := range txs {
		k := tx.OccurredAt.In(loc).Format("2006-01-02")
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
		day, err := time.ParseInLocation("2006-01-02", k, loc)
		if err != nil {
			continue
		}
		out = append(out, DayGroup{
			Label: FormatDayLabel(day, loc),
			Items: byDay[k],
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
