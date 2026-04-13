package handlers

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"moana/internal/store"
	"moana/internal/timeutil"
)

// History lists transactions with filters, search, sort, and date grouping.
func (a *App) History(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	loc := timeutil.LoadLocation(u.Timezone)
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	kindParam := strings.TrimSpace(r.URL.Query().Get("kind"))
	// Default: show all transactions. Use kind=income or kind=expense to narrow.
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
		// "", "all", or unknown → all
		kind = "all"
		filterKind = ""
	}
	sortParam := strings.TrimSpace(r.URL.Query().Get("sort"))
	oldestFirst := sortParam == "oldest"
	sortLabel := "newest"
	if oldestFirst {
		sortLabel = "oldest"
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	var f store.TransactionFilter
	f.Kind = filterKind
	f.Search = q
	f.OldestFirst = oldestFirst

	filterActive := from != "" && to != ""
	if filterActive {
		fu, tu, err := timeutil.DayRangeUTCFromLocalDates(loc, from, to)
		if err != nil {
			historyReturn := r.URL.RequestURI()
			if historyReturn == "" {
				historyReturn = "/history"
			}
			a.historyRender(w, u, historyPageData{
				Error:            "Invalid date range.",
				UserTZ:           u.Timezone,
				Kind:             kind,
				Search:           q,
				Sort:             sortLabel,
				FilterFrom:       from,
				FilterTo:         to,
				FilterActive:     true,
				Nav:              buildHistoryNav(r.URL),
				Groups:           nil,
				HistoryReturnURL: historyReturn,
			})
			return
		}
		f.FromUTC = &fu
		f.ToUTC = &tu
	}

	txs, err := a.Store.ListTransactions(ctx, u.ID, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	groups := groupTransactionsByDay(txs, loc, !oldestFirst)
	historyReturn := r.URL.RequestURI()
	if historyReturn == "" {
		historyReturn = "/history"
	}
	a.historyRender(w, u, historyPageData{
		Error:            "",
		UserTZ:           u.Timezone,
		Kind:             kind,
		Search:           q,
		Sort:             sortLabel,
		FilterFrom:       from,
		FilterTo:         to,
		FilterActive:     filterActive,
		Nav:              buildHistoryNav(r.URL),
		Groups:           groups,
		HistoryReturnURL: historyReturn,
	})
}

type historyPageData struct {
	Error            string
	UserTZ           string
	Kind             string
	Search           string
	Sort             string
	FilterFrom       string
	FilterTo         string
	FilterActive     bool
	Nav              historyNav
	Groups           []historyDayGroup
	HistoryReturnURL string // current /history path+query for edit "next" links
}

type historyNav struct {
	LinkAll     string
	LinkIncome  string
	LinkExpense string
	SortNewest  string
	SortOldest  string
}

type historyDayGroup struct {
	Label string
	Items []store.Transaction
}

func (a *App) historyRender(w http.ResponseWriter, u *store.User, data historyPageData) {
	ld := LayoutData{
		Title:  "History",
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: "history",
	}
	a.renderShell(w, "history_inner.html", data, ld)
}

func buildHistoryNav(u *url.URL) historyNav {
	with := func(mut func(v url.Values)) string {
		v := cloneQuery(u)
		mut(v)
		enc := v.Encode()
		if enc == "" {
			return "/history"
		}
		return "/history?" + enc
	}
	return historyNav{
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
	v := url.Values{}
	for k, vals := range u.Query() {
		for _, x := range vals {
			v.Add(k, x)
		}
	}
	return v
}

func groupTransactionsByDay(txs []store.Transaction, loc *time.Location, newestDayFirst bool) []historyDayGroup {
	if len(txs) == 0 {
		return nil
	}
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
	out := make([]historyDayGroup, 0, len(keys))
	for _, k := range keys {
		day, err := time.ParseInLocation("2006-01-02", k, loc)
		if err != nil {
			continue
		}
		out = append(out, historyDayGroup{
			Label: formatHistoryDayLabel(day, loc),
			Items: byDay[k],
		})
	}
	return out
}

func formatHistoryDayLabel(day time.Time, loc *time.Location) string {
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
