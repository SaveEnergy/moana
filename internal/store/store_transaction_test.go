package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestGetAndUpdateTransaction(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "upd@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, uid)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hid := owner.HouseholdID
	day := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uid, hid, -500, day, "a", nil)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := st.GetTransactionByID(ctx, hid, tid)
	if err != nil || tx == nil || tx.AmountCents != -500 {
		t.Fatalf("get: %+v err=%v", tx, err)
	}
	nilTx, err := st.GetTransactionByID(ctx, hid, 99999)
	if err != nil || nilTx != nil {
		t.Fatalf("missing row: %+v err=%v", nilTx, err)
	}
	newDay := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)
	if err := st.UpdateTransaction(ctx, hid, uid, tid, -600, newDay, "b", nil); err != nil {
		t.Fatal(err)
	}
	tx2, err := st.GetTransactionByID(ctx, hid, tid)
	if err != nil || tx2.Description != "b" || tx2.AmountCents != -600 {
		t.Fatalf("after: %+v", tx2)
	}
}
