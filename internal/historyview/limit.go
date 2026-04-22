package historyview

import "moana/internal/store"

// applyHistoryFetchLimit sets f.Limit to pageSize+1 so [trimHistoryRows] can detect "more rows exist" for the same filters.
func applyHistoryFetchLimit(f *store.TransactionFilter, pageSize int) (probe bool) {
	if f == nil {
		return false
	}
	if pageSize < 1 {
		pageSize = defaultHistoryPageSize
	}
	f.Limit = pageSize + 1
	return true
}

func trimHistoryRows(txs []store.Transaction, probe bool, pageSize int) ([]store.Transaction, bool) {
	if !probe || len(txs) <= pageSize {
		return txs, false
	}
	return txs[:pageSize], true
}
