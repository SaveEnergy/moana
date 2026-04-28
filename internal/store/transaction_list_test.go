package store

import (
	"testing"
	"time"
)

func TestListTransactions_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_recentLimitOnly_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{Limit: 3})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_kindIncomeRecentLimit_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{Kind: "income", Limit: 5})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_searchNoDate_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{Search: "x", Limit: 10})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_datedSingleBoundSearch_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	from := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{FromUTC: &from, Search: "q", Limit: 5})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_datedBothSearch_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{FromUTC: &from, ToUTC: &to, Search: "x", Limit: 10})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_datedNoSearch_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{FromUTC: &from, ToUTC: &to, Limit: 10})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_oldestFirstLimitOnly_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{Limit: 3, OldestFirst: true})
	assertErrIsContextCanceled(t, err)
}
