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
}
