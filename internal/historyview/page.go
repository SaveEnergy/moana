package historyview

import (
	"context"
	"net/url"
	"time"

	"moana/internal/store"
	"moana/internal/timeutil"
)

// BuildPage loads transactions and builds groups + nav for the history UI.
// requestURI should be r.URL.RequestURI() (or "" to default to /history).
func BuildPage(ctx context.Context, st *store.Store, householdID int64, loc *time.Location, u *url.URL, requestURI string) (PageData, error) {
	p := ParseHistoryURL(u)

	var f store.TransactionFilter
	f.Kind = p.filterKind
	f.Search = p.search
	f.OldestFirst = p.oldestFirst

	historyReturn := historyReturnOrDefault(requestURI)

	if partialDateFilter(p.from, p.to) {
		return PageData{
			Error:            "Invalid date range.",
			Kind:             p.kind,
			Search:           p.search,
			Sort:             p.sortLabel,
			FilterFrom:       p.from,
			FilterTo:         p.to,
			FilterActive:     true,
			Nav:              BuildNav(u),
			Groups:           nil,
			HistoryReturnURL: historyReturn,
			TruncationLimit:  defaultHistoryFetchLimit,
		}, nil
	}

	if p.filterActive {
		fu, tu, err := timeutil.DayRangeUTCFromLocalDates(loc, p.from, p.to)
		if err != nil {
			return PageData{
				Error:            "Invalid date range.",
				Kind:             p.kind,
				Search:           p.search,
				Sort:             p.sortLabel,
				FilterFrom:       p.from,
				FilterTo:         p.to,
				FilterActive:     true,
				Nav:              BuildNav(u),
				Groups:           nil,
				HistoryReturnURL: historyReturn,
				TruncationLimit:  defaultHistoryFetchLimit,
			}, nil
		}
		f.FromUTC = &fu
		f.ToUTC = &tu
	}

	probe := applyHistoryFetchLimit(&f)

	txs, err := st.ListTransactions(ctx, householdID, f)
	if err != nil {
		return PageData{}, err
	}
	txs, truncated := trimHistoryRows(txs, probe)
	groups := GroupByDay(txs, loc, !p.oldestFirst)
	return PageData{
		Error:            "",
		Kind:             p.kind,
		Search:           p.search,
		Sort:             p.sortLabel,
		FilterFrom:       p.from,
		FilterTo:         p.to,
		FilterActive:     p.filterActive,
		Nav:              BuildNav(u),
		Groups:           groups,
		HistoryReturnURL: historyReturn,
		Truncated:        truncated,
		TruncationLimit:  defaultHistoryFetchLimit,
	}, nil
}

// partialDateFilter is true when exactly one of from/to is set (both are required for filtering).
func partialDateFilter(from, to string) bool {
	return (from != "" || to != "") && (from == "" || to == "")
}
