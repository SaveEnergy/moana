package historyview

// RoutePath is the history ledger URL path. Must match [handlers.HistoryPath] and [safepath.Default].
const RoutePath = "/history"

// Query parameter names for history filter/sort links and [ParseHistoryURL].
const (
	QuerySearch = "q"
	QueryKind   = "kind"
	QuerySort   = "sort"
	QueryFrom   = "from"
	QueryTo     = "to"
	// QueryRows (n) caps how many transactions load per request (log-style; see historyview [defaultHistoryPageSize]).
	QueryRows = "n"
)

// SortOldestValue is the query value that means oldest-first ordering.
const SortOldestValue = "oldest"

// KindAll, KindIncome, and KindExpense are [QueryKind] values and normalized kind strings
// in PageData / [store.TransactionFilter].Kind ("income"/"expense" or "" for all kinds).
const (
	KindAll     = "all"
	KindIncome  = "income"
	KindExpense = "expense"
)

// SortLabelNewest is the display label for newest-first ordering ([PageData.Sort]).
// Oldest-first uses [SortOldestValue] as both query value and label.
const SortLabelNewest = "newest"
