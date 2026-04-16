package historyview

import (
	"maps"
	"net/url"
)

// BuildNav preserves the current query string while swapping kind/sort presets.
func BuildNav(u *url.URL) Nav {
	with := func(mut func(v url.Values)) string {
		v := cloneQuery(u)
		mut(v)
		enc := v.Encode()
		if enc == "" {
			return "/history"
		}
		return "/history?" + enc
	}
	return Nav{
		LinkAll: with(func(v url.Values) {
			v.Set("kind", "all")
		}),
		LinkIncome: with(func(v url.Values) {
			v.Set("kind", "income")
		}),
		LinkExpense: with(func(v url.Values) {
			v.Set("kind", "expense")
		}),
		SortNewest: with(func(v url.Values) {
			v.Del("sort")
		}),
		SortOldest: with(func(v url.Values) {
			v.Set("sort", "oldest")
		}),
	}
}

func cloneQuery(u *url.URL) url.Values {
	return maps.Clone(u.Query())
}
