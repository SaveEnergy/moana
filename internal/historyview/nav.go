package historyview

import (
	"maps"
	"net/url"
)

// BuildNav preserves the current query string while swapping kind/sort presets.
func BuildNav(u *url.URL) Nav {
	return buildNavFromValues(u.Query())
}

// buildNavFromValues builds nav links from url.Values (usually the same URL.Query map as parseHistoryURLValues).
func buildNavFromValues(q url.Values) Nav {
	base := maps.Clone(q)
	with := func(mut func(v url.Values)) string {
		v := maps.Clone(base)
		mut(v)
		return pathWithQuery(v.Encode())
	}
	return Nav{
		LinkAll: with(func(v url.Values) {
			v.Set(QueryKind, KindAll)
		}),
		LinkIncome: with(func(v url.Values) {
			v.Set(QueryKind, KindIncome)
		}),
		LinkExpense: with(func(v url.Values) {
			v.Set(QueryKind, KindExpense)
		}),
		SortNewest: with(func(v url.Values) {
			v.Del(QuerySort)
		}),
		SortOldest: with(func(v url.Values) {
			v.Set(QuerySort, SortOldestValue)
		}),
	}
}

// pathWithQuery returns [RoutePath] alone or with a non-empty encoded query (?…).
func pathWithQuery(encodedQuery string) string {
	if encodedQuery == "" {
		return RoutePath
	}
	return RoutePath + "?" + encodedQuery
}
