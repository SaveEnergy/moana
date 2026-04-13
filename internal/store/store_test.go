package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/auth"
	"moana/internal/db"
)

func TestUserCategoryTransactionFlow(t *testing.T) {
	t.Parallel()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	st := New(database)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-test-123")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "flow@example.com", hash, "user", "Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "flow@example.com")
	if err != nil || u == nil || u.ID != uid {
		t.Fatalf("user: %+v err=%v", u, err)
	}

	catID, err := st.CreateCategory(ctx, uid, "Salary")
	if err != nil {
		t.Fatal(err)
	}
	cats, err := st.ListCategories(ctx, uid)
	if err != nil || len(cats) != 1 {
		t.Fatalf("categories: %v err=%v", cats, err)
	}

	day := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uid, 50000, day, "paycheck", &catID)
	if err != nil {
		t.Fatal(err)
	}
	if tid <= 0 {
		t.Fatal("expected transaction id")
	}

	sum, err := st.SumAmountCents(ctx, uid, nil, nil)
	if err != nil || sum != 50000 {
		t.Fatalf("sum=%d err=%v", sum, err)
	}

	txs, err := st.ListTransactions(ctx, uid, TransactionFilter{})
	if err != nil || len(txs) != 1 {
		t.Fatalf("list: %v err=%v", txs, err)
	}
	if txs[0].AmountCents != 50000 || txs[0].CategoryName != "Salary" {
		t.Fatalf("row: %+v", txs[0])
	}

	if err := st.DeleteCategory(ctx, uid, catID); err != nil {
		t.Fatal(err)
	}
	txs2, err := st.ListTransactions(ctx, uid, TransactionFilter{})
	if err != nil || len(txs2) != 1 {
		t.Fatal(err)
	}
	if txs2[0].CategoryID.Valid {
		t.Fatal("expected category cleared")
	}
}

func TestGetAndUpdateTransaction(t *testing.T) {
	t.Parallel()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	st := New(database)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "upd@example.com", hash, "user", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uid, -500, day, "a", nil)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := st.GetTransactionByID(ctx, uid, tid)
	if err != nil || tx == nil || tx.AmountCents != -500 {
		t.Fatalf("get: %+v err=%v", tx, err)
	}
	nilTx, err := st.GetTransactionByID(ctx, uid, 99999)
	if err != nil || nilTx != nil {
		t.Fatalf("missing row: %+v err=%v", nilTx, err)
	}
	newDay := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)
	if err := st.UpdateTransaction(ctx, uid, tid, -600, newDay, "b", nil); err != nil {
		t.Fatal(err)
	}
	tx2, err := st.GetTransactionByID(ctx, uid, tid)
	if err != nil || tx2.Description != "b" || tx2.AmountCents != -600 {
		t.Fatalf("after: %+v", tx2)
	}
}

func TestListUsers(t *testing.T) {
	t.Parallel()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	st := New(database)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "a@example.com", hash, "user", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "b@example.com", hash, "admin", "Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	users, err := st.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("len=%d", len(users))
	}
	if users[0].Email != "a@example.com" || users[1].Role != "admin" {
		t.Fatalf("rows: %+v %+v", users[0], users[1])
	}
}
