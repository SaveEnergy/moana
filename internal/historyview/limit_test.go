package historyview

import (
	"testing"

	"moana/internal/store"
)

func TestApplyHistoryFetchLimit(t *testing.T) {
	t.Parallel()
	t.Run("sets default when unset", func(t *testing.T) {
		t.Parallel()
		f := store.TransactionFilter{}
		if probe := applyHistoryFetchLimit(&f); !probe {
			t.Fatal("want probe when applying default cap")
		}
		want := defaultHistoryFetchLimit + 1
		if f.Limit != want {
			t.Fatalf("Limit=%d want %d (cap+1 probe)", f.Limit, want)
		}
	})
	t.Run("preserves explicit limit", func(t *testing.T) {
		t.Parallel()
		f := store.TransactionFilter{Limit: 42}
		if probe := applyHistoryFetchLimit(&f); probe {
			t.Fatal("unexpected probe")
		}
		if f.Limit != 42 {
			t.Fatalf("Limit=%d want 42", f.Limit)
		}
	})
	t.Run("nil filter is noop", func(t *testing.T) {
		t.Parallel()
		applyHistoryFetchLimit(nil) // must not panic
	})
}

func TestTrimHistoryRows(t *testing.T) {
	t.Parallel()
	makeTxs := func(n int) []store.Transaction {
		out := make([]store.Transaction, n)
		for i := range out {
			out[i].ID = int64(i + 1)
		}
		return out
	}
	t.Run("no probe never truncates", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(defaultHistoryFetchLimit + 10)
		out, tr := trimHistoryRows(txs, false)
		if tr || len(out) != len(txs) {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
	t.Run("probe under cap", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(10)
		out, tr := trimHistoryRows(txs, true)
		if tr || len(out) != 10 {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
	t.Run("probe at cap not truncated", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(defaultHistoryFetchLimit)
		out, tr := trimHistoryRows(txs, true)
		if tr || len(out) != defaultHistoryFetchLimit {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
	t.Run("probe over cap", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(defaultHistoryFetchLimit + 1)
		out, tr := trimHistoryRows(txs, true)
		if !tr || len(out) != defaultHistoryFetchLimit {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
}
