package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestListUsers(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "b@example.com", hash, "admin")
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

func TestHouseholdMembersSeeSameTransactions(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-a")
	if err != nil {
		t.Fatal(err)
	}
	ownerID, err := st.CreateUser(ctx, "owner@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, ownerID)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hid := owner.HouseholdID

	hashB, err := auth.HashPassword("pw-b")
	if err != nil {
		t.Fatal(err)
	}
	memberID, err := st.CreateHouseholdMember(ctx, hid, "member@example.com", hashB)
	if err != nil {
		t.Fatal(err)
	}
	catID, err := st.CreateCategory(ctx, hid, "Food", "", "")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, ownerID, hid, -999, day, "shared lunch", &catID)
	if err != nil {
		t.Fatal(err)
	}

	listOwner, err := st.ListTransactions(ctx, hid, TransactionFilter{})
	if err != nil || len(listOwner) != 1 || listOwner[0].ID != tid {
		t.Fatalf("owner list: %+v err=%v", listOwner, err)
	}
	member, err := st.GetUserByID(ctx, memberID)
	if err != nil || member == nil || member.HouseholdID != hid {
		t.Fatalf("member: %+v", member)
	}
	listMember, err := st.ListTransactions(ctx, member.HouseholdID, TransactionFilter{})
	if err != nil || len(listMember) != 1 || listMember[0].Description != "shared lunch" {
		t.Fatalf("member list: %+v err=%v", listMember, err)
	}
	got, err := st.GetTransactionByID(ctx, hid, tid)
	if err != nil || got == nil || got.UserID != ownerID {
		t.Fatalf("get by household: %+v", got)
	}
}

func TestCreateUser_duplicateEmail(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateUser(ctx, "same@example.com", hash, "user"); err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "same@example.com", hash, "user")
	if !errors.Is(err, ErrDuplicateUserEmail) {
		t.Fatalf("got %v want %v", err, ErrDuplicateUserEmail)
	}
}

func TestCreateHouseholdMember_duplicateEmail(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "dup-owner@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, uid)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hid := owner.HouseholdID
	hash2, err := auth.HashPassword("y")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateHouseholdMember(ctx, hid, "dup-mem@example.com", hash2); err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateHouseholdMember(ctx, hid, "dup-mem@example.com", hash2)
	if !errors.Is(err, ErrDuplicateUserEmail) {
		t.Fatalf("got %v want %v", err, ErrDuplicateUserEmail)
	}
}
