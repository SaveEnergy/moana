package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/passwordtest"
)

func TestGetAndUpdateTransaction(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
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

func TestGetTransactionByID_wrongHouseholdReturnsNil(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uidA, err := st.CreateUser(ctx, "gtx-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uA, err := st.GetUserByID(ctx, uidA)
	if err != nil || uA == nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uidA, uA.HouseholdID, -100, day, "solo", nil)
	if err != nil {
		t.Fatal(err)
	}

	uidB, err := st.CreateUser(ctx, "gtx-b@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uB, err := st.GetUserByID(ctx, uidB)
	if err != nil || uB == nil {
		t.Fatal(err)
	}
	if uA.HouseholdID == uB.HouseholdID {
		t.Fatal("expected distinct households")
	}

	got, err := st.GetTransactionByID(ctx, uB.HouseholdID, tid)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("transaction from another household must not resolve: %+v", got)
	}
}

func TestGetTransactionByID_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.GetTransactionByID(alreadyCancelledContext(t), 1, 1)
	assertErrIsContextCanceled(t, err)
}
