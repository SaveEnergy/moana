package historyview

import (
	"testing"

	"moana/internal/store"
)

func TestApplyHistoryFetchLimit(t *testing.T) {
	t.Parallel()
	t.Run("sets limit to pageSize+1", func(t *testing.T) {
		t.Parallel()
		f := store.TransactionFilter{}
		if probe := applyHistoryFetchLimit(&f, defaultHistoryPageSize); !probe {
			t.Fatal("want probe when applying cap")
		}
		want := defaultHistoryPageSize + 1
		if f.Limit != want {
			t.Fatalf("Limit=%d want %d (cap+1 probe)", f.Limit, want)
		}
	})
	t.Run("nil filter is noop", func(t *testing.T) {
		t.Parallel()
		applyHistoryFetchLimit(nil, defaultHistoryPageSize) // must not panic
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
	ps := 100
	t.Run("no probe never truncates", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(ps + 10)
		out, tr := trimHistoryRows(txs, false, ps)
		if tr || len(out) != len(txs) {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
	t.Run("probe under cap", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(10)
		out, tr := trimHistoryRows(txs, true, ps)
		if tr || len(out) != 10 {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
	t.Run("probe at cap not truncated", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(ps)
		out, tr := trimHistoryRows(txs, true, ps)
		if tr || len(out) != ps {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
	t.Run("probe over cap", func(t *testing.T) {
		t.Parallel()
		txs := makeTxs(ps + 1)
		out, tr := trimHistoryRows(txs, true, ps)
		if !tr || len(out) != ps {
			t.Fatalf("len=%d truncated=%v", len(out), tr)
		}
	})
}
