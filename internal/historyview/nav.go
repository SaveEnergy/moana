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
		enc := v.Encode()
		if enc == "" {
			return RoutePath
		}
		return RoutePath + "?" + enc
	}
	return Nav{
		LinkAll: with(func(v url.Values) {
			v.Set(QueryKind, "all")
		}),
		LinkIncome: with(func(v url.Values) {
			v.Set(QueryKind, "income")
		}),
		LinkExpense: with(func(v url.Values) {
			v.Set(QueryKind, "expense")
		}),
		SortNewest: with(func(v url.Values) {
			v.Del(QuerySort)
		}),
		SortOldest: with(func(v url.Values) {
			v.Set(QuerySort, SortOldestValue)
		}),
	}
}
