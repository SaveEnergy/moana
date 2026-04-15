package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestCreateTransaction_rejectsOtherHouseholdCategory(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-mutate-test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "mutate-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uA, err := st.GetUserByEmail(ctx, "mutate-a@example.com")
	if err != nil || uA == nil {
		t.Fatal(err)
	}
	catID, err := st.CreateCategory(ctx, uA.HouseholdID, "Food", "", "")
	if err != nil {
		t.Fatal(err)
	}

	uidB, err := st.CreateUser(ctx, "mutate-b@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uB, err := st.GetUserByEmail(ctx, "mutate-b@example.com")
	if err != nil || uB == nil {
		t.Fatal(err)
	}
	if uB.HouseholdID == uA.HouseholdID {
		t.Fatal("expected distinct households")
	}

	day := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)
	_, err = st.CreateTransaction(ctx, uidB, uB.HouseholdID, -500, day, "x", &catID)
	if err == nil {
		t.Fatal("expected ErrInvalidCategory")
	}
	if !errors.Is(err, ErrInvalidCategory) {
		t.Fatalf("got %v want %v", err, ErrInvalidCategory)
	}
}

func TestUpdateTransaction_rejectsOtherHouseholdCategory(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-mutate-up")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "mutate-up-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uA, err := st.GetUserByEmail(ctx, "mutate-up-a@example.com")
	if err != nil || uA == nil {
		t.Fatal(err)
	}
	catA, err := st.CreateCategory(ctx, uA.HouseholdID, "Rent", "", "")
	if err != nil {
		t.Fatal(err)
	}

	uidB, err := st.CreateUser(ctx, "mutate-up-b@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uB, err := st.GetUserByEmail(ctx, "mutate-up-b@example.com")
	if err != nil || uB == nil {
		t.Fatal(err)
	}

	day := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uidB, uB.HouseholdID, -100, day, "solo", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = st.UpdateTransaction(ctx, uB.HouseholdID, uidB, tid, -200, day, "x", &catA)
	if err == nil {
		t.Fatal("expected ErrInvalidCategory")
	}
	if !errors.Is(err, ErrInvalidCategory) {
		t.Fatalf("got %v", err)
	}
}
