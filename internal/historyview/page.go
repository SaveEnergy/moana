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
	q := u.Query()
	p := parseHistoryURLValues(q)

	var f store.TransactionFilter
	f.Kind = p.filterKind
	f.Search = p.search
	f.OldestFirst = p.oldestFirst

	historyReturn := historyReturnOrDefault(requestURI)

	if partialDateFilter(p.from, p.to) {
		return invalidDateRangePage(p, historyReturn, q, loc), nil
	}

	if p.filterActive {
		fu, tu, err := timeutil.DayRangeUTCFromLocalDates(loc, p.from, p.to)
		if err != nil {
			return invalidDateRangePage(p, historyReturn, q, loc), nil
		}
		f.FromUTC = &fu
		f.ToUTC = &tu
	}

	pageSize := p.rows
	probe := applyHistoryFetchLimit(&f, pageSize)

	txs, err := st.ListTransactions(ctx, householdID, f)
	if err != nil {
		return PageData{}, err
	}
	txs, truncated := trimHistoryRows(txs, probe, pageSize)
	groups := GroupByDay(txs, loc, !p.oldestFirst)
	q7, q30, q90 := buildAllQuickRangeLinks(q, loc)
	return PageData{
		Error:            "",
		Kind:             p.kind,
		Search:           p.search,
		Sort:             p.sortLabel,
		Rows:             pageSize,
		RowOptions:       DefaultHistoryPageSizes,
		QuickLast7URL:    q7,
		QuickLast30URL:   q30,
		QuickLast90URL:   q90,
		FilterFrom:       p.from,
		FilterTo:         p.to,
		FilterActive:     p.filterActive,
		Nav:              buildNavFromValues(q),
		Groups:           groups,
		HistoryReturnURL: historyReturn,
		Truncated:        truncated,
		TruncationLimit:  pageSize,
	}, nil
}

// invalidDateRangePage builds the standard /history payload when from/to cannot be applied.
func invalidDateRangePage(p HistoryURLParams, historyReturn string, q url.Values, loc *time.Location) PageData {
	q7, q30, q90 := buildAllQuickRangeLinks(q, loc)
	return PageData{
		Error:            InvalidDateRangeMessage,
		Kind:             p.kind,
		Search:           p.search,
		Sort:             p.sortLabel,
		Rows:             p.rows,
		RowOptions:       DefaultHistoryPageSizes,
		QuickLast7URL:    q7,
		QuickLast30URL:   q30,
		QuickLast90URL:   q90,
		FilterFrom:       p.from,
		FilterTo:         p.to,
		FilterActive:     true,
		Nav:              buildNavFromValues(q),
		Groups:           nil,
		HistoryReturnURL: historyReturn,
		TruncationLimit:  p.rows,
	}
}

// partialDateFilter is true when exactly one of from/to is set (both are required for filtering).
func partialDateFilter(from, to string) bool {
	return (from != "" || to != "") && (from == "" || to == "")
}
