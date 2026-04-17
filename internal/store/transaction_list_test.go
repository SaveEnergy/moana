package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestListTransactions_respectsLimit(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-limit-test")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "limit-list@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "limit-list@example.com")
	if err != nil || u == nil {
		t.Fatalf("user: %+v err=%v", u, err)
	}
	hid := u.HouseholdID

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

func TestListTransactions_searchLiteralPercentDoesNotMatchAllRows(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "like-pct@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -200, day, "a%b", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Search: "%"})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].Description != "a%b" {
		t.Fatalf("want 1 row with literal %%, got %+v", txs)
	}
}

func TestListTransactions_whitespaceOnlySearchSkipped(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "ws-search@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Search: " \t  "})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 {
		t.Fatalf("whitespace-only search should not filter, got %d rows", len(txs))
	}
}

func TestListTransactions_kindIgnoresSurroundingSpace(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "kind-trim@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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
