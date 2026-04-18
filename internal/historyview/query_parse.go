package historyview

import (
	"net/url"
	"strings"
)

// HistoryURLParams holds normalized /history query string fields.
type HistoryURLParams struct {
	kind         string
	filterKind   string
	search       string
	sortLabel    string
	oldestFirst  bool
	from         string
	to           string
	filterActive bool
}

// ParseHistoryURL extracts normalized filters from a /history URL.
func ParseHistoryURL(u *url.URL) HistoryURLParams {
	return parseHistoryURLValues(u.Query())
}

// parseHistoryURLValues normalizes filters from a url.Values produced by URL.Query.
// BuildPage passes the same Values to this and to buildNavFromValues (one query parse per request).
func parseHistoryURLValues(v url.Values) HistoryURLParams {
	q := strings.TrimSpace(v.Get("q"))
	kindParam := strings.TrimSpace(v.Get("kind"))
	kind := "all"
	filterKind := ""
	switch kindParam {
	case "income":
		kind = "income"
		filterKind = "income"
	case "expense":
		kind = "expense"
		filterKind = "expense"
	default:
		kind = "all"
		filterKind = ""
	}
	sortParam := strings.TrimSpace(v.Get("sort"))
	oldestFirst := sortParam == "oldest"
	sortLabel := "newest"
	if oldestFirst {
		sortLabel = "oldest"
	}
	from := strings.TrimSpace(v.Get("from"))
	to := strings.TrimSpace(v.Get("to"))
	filterActive := from != "" && to != ""
	return HistoryURLParams{
		kind:         kind,
		filterKind:   filterKind,
		search:       q,
		sortLabel:    sortLabel,
		oldestFirst:  oldestFirst,
		from:         from,
		to:           to,
		filterActive: filterActive,
	}
}

func historyReturnOrDefault(requestURI string) string {
	if requestURI == "" {
		return "/history"
	}
	return requestURI
}
