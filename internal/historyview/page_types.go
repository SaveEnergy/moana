package historyview

import "moana/internal/store"

// InvalidDateRangeMessage is shown when [BuildPage] cannot apply date filters (partial from/to or unparseable dates).
const InvalidDateRangeMessage = "Invalid date range."

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
	Error  string
	Kind   string
	Search string
	Sort   string
	// Rows is the active n= cap (or default) — template fills the row-limit select.
	Rows       int
	RowOptions []int
	// Quick* are GET /history links that set from/to to inclusive local day ranges, preserving the rest of the query.
	QuickLast7URL  string
	QuickLast30URL string
	QuickLast90URL string
	FilterFrom     string
	FilterTo       string
	FilterActive   bool
	Nav            Nav
	Groups         []DayGroup
	HistoryReturnURL string // current /history path+query for edit "next" links
	// Truncated is true when more rows matched than the active row cap.
	Truncated       bool
	TruncationLimit int // same as Rows in normal builds (invalid-date payload uses parsed n= or default)
}
