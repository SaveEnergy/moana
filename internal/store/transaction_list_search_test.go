package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/passwordtest"
)

func TestListTransactions_searchNoDate_incomeMatchesDescription(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "list-search-nodate@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 10000, day, "windfall bonus", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -3000, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "income", Search: "bonus", Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].Description != "windfall bonus" {
		t.Fatalf("want 1 row, got %+v", txs)
	}
}

func TestListTransactions_searchNoDate_zeroLimit_returnsAllMatches(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "list-search-z@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, day, "alpha token", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 200, day, "beta token", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Search: "token", Limit: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 2 {
		t.Fatalf("want 2 rows, got %d %+v", len(txs), txs)
	}
}

func TestListTransactions_datedSingleBoundSearch_fromOnlyExpense(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "list-1bnd@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	from := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	day := time.Date(2026, 12, 10, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -500, day, "matchme latte", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 10000, day, "salary noise", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{FromUTC: &from, Search: "matchme", Kind: "expense", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].Description != "matchme latte" {
		t.Fatalf("want 1 row, got %+v", txs)
	}
}

func TestListTransactions_searchLiteralPercentDoesNotMatchAllRows(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
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
	hash := passwordtest.MustHash(t, "x")
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

func TestListTransactions_combinedFilterOldestFirstWithSearchAndBounds(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "combo-list@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	day1 := time.Date(2026, 8, 10, 15, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 8, 11, 15, 0, 0, 0, time.UTC)
	from := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 10, 23, 59, 59, 0, time.UTC)

	if _, err := st.CreateTransaction(ctx, uid, hid, 500, day1, "salary august", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, day1, "coffee shop", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 300, day2, "salary next day", nil); err != nil {
		t.Fatal(err)
	}

	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{
		FromUTC:     &from,
		ToUTC:       &to,
		Kind:        "income",
		Search:      "sal",
		OldestFirst: true,
		Limit:       10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].Description != "salary august" || txs[0].AmountCents != 500 {
		t.Fatalf("want 1 income row on day1 matching search, got %+v", txs)
	}
}
