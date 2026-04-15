package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestUserCategoryTransactionFlow(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-test-123")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "flow@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "flow@example.com")
	if err != nil || u == nil || u.ID != uid {
		t.Fatalf("user: %+v err=%v", u, err)
	}
	hid := u.HouseholdID

	catID, err := st.CreateCategory(ctx, hid, "Salary", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cats, err := st.ListCategories(ctx, hid)
	if err != nil || len(cats) != 1 {
		t.Fatalf("categories: %v err=%v", cats, err)
	}

	day := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uid, hid, 50000, day, "paycheck", &catID)
	if err != nil {
		t.Fatal(err)
	}
	if tid <= 0 {
		t.Fatal("expected transaction id")
	}

	sum, err := st.SumAmountCents(ctx, hid, nil, nil)
	if err != nil || sum != 50000 {
		t.Fatalf("sum=%d err=%v", sum, err)
	}

	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{})
	if err != nil || len(txs) != 1 {
		t.Fatalf("list: %v err=%v", txs, err)
	}
	if txs[0].AmountCents != 50000 || txs[0].CategoryName != "Salary" {
		t.Fatalf("row: %+v", txs[0])
	}

	if err := st.DeleteCategory(ctx, uid, catID); err != nil {
		t.Fatal(err)
	}
	txs2, err := st.ListTransactions(ctx, hid, TransactionFilter{})
	if err != nil || len(txs2) != 1 {
		t.Fatal(err)
	}
	if txs2[0].CategoryID.Valid {
		t.Fatal("expected category cleared")
	}
}
