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
)

// SortOldestValue is the query value that means oldest-first ordering.
const SortOldestValue = "oldest"
