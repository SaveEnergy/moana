package historyview

import "moana/internal/store"

// defaultHistoryFetchLimit caps how many transactions the history page loads in one request.
// Without a cap, ListTransactions would scan the full household ledger into memory.
const defaultHistoryFetchLimit = 5000

// applyHistoryFetchLimit sets a safe row cap when none was requested. It fetches one extra row
// so BuildPage can tell whether results were truncated (see trimHistoryRows).
func applyHistoryFetchLimit(f *store.TransactionFilter) (probeTruncation bool) {
	if f == nil || f.Limit > 0 {
		return false
	}
	f.Limit = defaultHistoryFetchLimit + 1
	return true
}

func trimHistoryRows(txs []store.Transaction, probe bool) ([]store.Transaction, bool) {
	if !probe || len(txs) <= defaultHistoryFetchLimit {
		return txs, false
	}
	return txs[:defaultHistoryFetchLimit], true
}
