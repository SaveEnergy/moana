package historyview

import "moana/internal/store"

// Nav holds sidebar link targets for the history filter/sort controls.
type Nav struct {
	LinkAll     string
	LinkIncome  string
	LinkExpense string
	SortNewest  string
	SortOldest  string
}

// DayGroup is one collapsible day section on the history page.
type DayGroup struct {
	Label string
	Items []store.Transaction
}

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
