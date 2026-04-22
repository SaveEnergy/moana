package historyview

import (
	"maps"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Default page size (newest-first cap) when /history is loaded with no n= param — avoids loading the full ledger.
const defaultHistoryPageSize = 100

// hardMaxHistoryPageSize is the largest allowed n= (higher values fall back to this).
const hardMaxHistoryPageSize = 2000

// DefaultHistoryPageSizes is the set of valid n= query values (log-style "tail" limits).
// Exported for [PageData.RowOptions] in templates.
var DefaultHistoryPageSizes = []int{50, 100, 250, 500, 2000}

// parseHistoryPageSize returns a normalized row cap from the n= query (default [defaultHistoryPageSize]).
func parseHistoryPageSize(v url.Values) int {
	s := strings.TrimSpace(v.Get(QueryRows))
	if s == "" {
		return defaultHistoryPageSize
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultHistoryPageSize
	}
	for _, a := range DefaultHistoryPageSizes {
		if n == a {
			return a
		}
	}
	if n > hardMaxHistoryPageSize {
		return hardMaxHistoryPageSize
	}
	// Unknown values snap to default so ad-hoc n=50000 cannot bypass the cap.
	return defaultHistoryPageSize
}

// buildQuickRangeLinks returns GET /history links that set from/to to inclusive local day ranges, preserving the rest of the query.
func buildQuickRangeLinks(q url.Values, loc *time.Location, days int) string {
	if loc == nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	toDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	fromDay := toDay.AddDate(0, 0, -(days - 1))
	v := maps.Clone(q)
	v.Set(QueryFrom, fromDay.Format("2006-01-02"))
	v.Set(QueryTo, toDay.Format("2006-01-02"))
	return pathWithQuery(v.Encode())
}

func buildAllQuickRangeLinks(q url.Values, loc *time.Location) (last7, last30, last90 string) {
	return buildQuickRangeLinks(q, loc, 7), buildQuickRangeLinks(q, loc, 30), buildQuickRangeLinks(q, loc, 90)
}
