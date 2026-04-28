package store

import (
	"context"
	"testing"
	"time"
)

func TestListTransactions_kindExpenseWithLimit_newestFirst(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "kind-exp-limit@example.com")
	base := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, base, "old", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -200, base.Add(time.Hour), "mid", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -300, base.Add(2*time.Hour), "new", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "expense", Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 2 {
		t.Fatalf("want 2, got %d %+v", len(txs), txs)
	}
	if txs[0].Description != "new" || txs[1].Description != "mid" {
		t.Fatalf("want newest first: %+v %+v", txs[0], txs[1])
	}
}

func TestListTransactions_oldestFirstIncomeWithLimit_oldestFirst(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "old-inc@example.com")
	base := time.Date(2026, 12, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, base, "first", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 200, base.Add(time.Hour), "second", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 300, base.Add(2*time.Hour), "third", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "income", Limit: 2, OldestFirst: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 2 {
		t.Fatalf("want 2 got %d %+v", len(txs), txs)
	}
	if txs[0].Description != "first" || txs[1].Description != "second" {
		t.Fatalf("want oldest first: %+v %+v", txs[0], txs[1])
	}
}

func TestListTransactions_respectsLimit(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "limit-list@example.com")

	base := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	for i := range 5 {
		occ := base.Add(time.Duration(i) * time.Hour)
		if _, err := st.CreateTransaction(ctx, uid, hid, int64(100+i), occ, "tx", nil); err != nil {
			t.Fatal(err)
		}
	}

	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 3 {
		t.Fatalf("want 3 rows, got %d", len(txs))
	}
	// Newest first (default): last inserted hour wins
	if txs[0].AmountCents != 104 || txs[2].AmountCents != 102 {
		t.Fatalf("unexpected order: %+v %+v", txs[0], txs[2])
	}

	txsAll, err := st.ListTransactions(ctx, hid, TransactionFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(txsAll) != 5 {
		t.Fatalf("unlimited list: want 5, got %d", len(txsAll))
	}

	txsZero, err := st.ListTransactions(ctx, hid, TransactionFilter{Limit: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(txsZero) != len(txsAll) {
		t.Fatalf("Limit 0: want %d rows like empty filter, got %d", len(txsAll), len(txsZero))
	}

	txsNeg, err := st.ListTransactions(ctx, hid, TransactionFilter{Limit: -1})
	if err != nil {
		t.Fatal(err)
	}
	if len(txsNeg) != 5 {
		t.Fatalf("negative limit treated as no cap: want 5, got %d", len(txsNeg))
	}
}

func TestListTransactions_kindIgnoresSurroundingSpace(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "kind-trim@example.com")
	day := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -50, day, "out", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "  income  "})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].AmountCents != 100 {
		t.Fatalf("want 1 income row, got %+v", txs)
	}
}

func TestListTransactions_kindExpenseOnlyNegativeAmounts(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "kind-expense@example.com")
	day := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -50, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "expense"})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].AmountCents != -50 || txs[0].Description != "coffee" {
		t.Fatalf("want 1 expense row, got %+v", txs)
	}
}

func TestListTransactions_unknownKindDoesNotFilterBySign(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "kind-unknown@example.com")
	day := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -50, day, "out", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "not-a-valid-kind"})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 2 {
		t.Fatalf("unknown kind must not apply income/expense SQL filter: want 2 rows, got %d", len(txs))
	}
}
