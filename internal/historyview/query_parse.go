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
	q := strings.TrimSpace(u.Query().Get("q"))
	kindParam := strings.TrimSpace(u.Query().Get("kind"))
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
	sortParam := strings.TrimSpace(u.Query().Get("sort"))
	oldestFirst := sortParam == "oldest"
	sortLabel := "newest"
	if oldestFirst {
		sortLabel = "oldest"
	}
	from := u.Query().Get("from")
	to := u.Query().Get("to")
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
