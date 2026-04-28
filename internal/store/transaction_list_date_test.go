package store

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestListTransactions_datedNoSearch_boundsIncomeOldestFirst(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "list-dated-ns@example.com")
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
		OldestFirst: true,
		Limit:       10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].Description != "salary august" || txs[0].AmountCents != 500 {
		t.Fatalf("want 1 income row on day1, got %+v", txs)
	}
}

func TestListTransactions_datedNoSearch_zeroLimit_returnsAllInRange(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "list-dated-z@example.com")
	from := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 3, 31, 23, 59, 59, 0, time.UTC)
	for i := 0; i < 3; i++ {
		d := time.Date(2026, 3, 5+i, 12, 0, 0, 0, time.UTC)
		if _, err := st.CreateTransaction(ctx, uid, hid, int64(-100*(i+1)), d, fmt.Sprintf("t%d", i), nil); err != nil {
			t.Fatal(err)
		}
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{FromUTC: &from, ToUTC: &to, Limit: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 3 {
		t.Fatalf("want 3 rows, got %d %+v", len(txs), txs)
	}
}
