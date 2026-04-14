package historyview

import (
	"context"
	"net/url"
	"strings"
	"time"

	"moana/internal/store"
	"moana/internal/timeutil"
)

// PageData is the template payload for the history ledger page.
type PageData struct {
	Error            string
	Kind             string
	Search           string
	Sort             string
	FilterFrom       string
	FilterTo         string
	FilterActive     bool
	Nav              Nav
	Groups           []DayGroup
	HistoryReturnURL string // current /history path+query for edit "next" links
}

// BuildPage loads transactions and builds groups + nav for the history UI.
// requestURI should be r.URL.RequestURI() (or "" to default to /history).
func BuildPage(ctx context.Context, st *store.Store, userID int64, loc *time.Location, u *url.URL, requestURI string) (PageData, error) {
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

	var f store.TransactionFilter
	f.Kind = filterKind
	f.Search = q
	f.OldestFirst = oldestFirst

	filterActive := from != "" && to != ""
	if filterActive {
		fu, tu, err := timeutil.DayRangeUTCFromLocalDates(loc, from, to)
		if err != nil {
			historyReturn := requestURI
			if historyReturn == "" {
				historyReturn = "/history"
			}
			return PageData{
				Error:            "Invalid date range.",
				Kind:             kind,
				Search:           q,
				Sort:             sortLabel,
				FilterFrom:       from,
				FilterTo:         to,
				FilterActive:     true,
				Nav:              BuildNav(u),
				Groups:           nil,
				HistoryReturnURL: historyReturn,
			}, nil
		}
		f.FromUTC = &fu
		f.ToUTC = &tu
	}

	txs, err := st.ListTransactions(ctx, userID, f)
	if err != nil {
		return PageData{}, err
	}
	groups := GroupByDay(txs, loc, !oldestFirst)
	historyReturn := requestURI
	if historyReturn == "" {
		historyReturn = "/history"
	}
	return PageData{
		Error:            "",
		Kind:             kind,
		Search:           q,
		Sort:             sortLabel,
		FilterFrom:       from,
		FilterTo:         to,
		FilterActive:     filterActive,
		Nav:              BuildNav(u),
		Groups:           groups,
		HistoryReturnURL: historyReturn,
	}, nil
}
