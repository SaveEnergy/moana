package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"moana/internal/passwordtest"
)

func TestCreateTransaction_rejectsOtherHouseholdCategory(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash := passwordtest.MustHash(t, "pw-mutate-test")
	_, err := st.CreateUser(ctx, "mutate-a@example.com", hash, "user")
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

func TestCreateTransaction_rejectsUserNotInHousehold(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "pw-mutate-hh-create")

	_, err := st.CreateUser(ctx, "cth-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uA, err := st.GetUserByEmail(ctx, "cth-a@example.com")
	if err != nil || uA == nil {
		t.Fatal(err)
	}
	uidB, err := st.CreateUser(ctx, "cth-b@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uB, err := st.GetUserByEmail(ctx, "cth-b@example.com")
	if err != nil || uB == nil {
		t.Fatal(err)
	}
	if uA.HouseholdID == uB.HouseholdID {
		t.Fatal("expected distinct households")
	}

	day := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	_, err = st.CreateTransaction(ctx, uidB, uA.HouseholdID, -100, day, "wrong household", nil)
	if err == nil {
		t.Fatal("expected ErrUserNotInHousehold")
	}
	if !errors.Is(err, ErrUserNotInHousehold) {
		t.Fatalf("got %v want %v", err, ErrUserNotInHousehold)
	}
}

func TestUpdateTransaction_rejectsOtherHouseholdCategory(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash := passwordtest.MustHash(t, "pw-mutate-up")
	_, err := st.CreateUser(ctx, "mutate-up-a@example.com", hash, "user")
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

func TestUpdateTransaction_wrongHouseholdReturnsNoRows(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "pw-mutate-hh")

	uidA, err := st.CreateUser(ctx, "mutate-hh-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uA, err := st.GetUserByEmail(ctx, "mutate-hh-a@example.com")
	if err != nil || uA == nil {
		t.Fatal(err)
	}

	uidB, err := st.CreateUser(ctx, "mutate-hh-b@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uB, err := st.GetUserByEmail(ctx, "mutate-hh-b@example.com")
	if err != nil || uB == nil {
		t.Fatal(err)
	}
	if uA.HouseholdID == uB.HouseholdID {
		t.Fatal("expected distinct households")
	}

	day := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uidA, uA.HouseholdID, -100, day, "mine", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = st.UpdateTransaction(ctx, uB.HouseholdID, uidB, tid, -200, day, "hijack", nil)
	if err == nil {
		t.Fatal("expected sql.ErrNoRows")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("got %v want %v", err, sql.ErrNoRows)
	}
}

func TestUpdateTransaction_actorNotInHouseholdReturnsNoRows(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "pw-mutate-actor")

	uidA, err := st.CreateUser(ctx, "mutate-act-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uA, err := st.GetUserByEmail(ctx, "mutate-act-a@example.com")
	if err != nil || uA == nil {
		t.Fatal(err)
	}
	uidB, err := st.CreateUser(ctx, "mutate-act-b@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uB, err := st.GetUserByEmail(ctx, "mutate-act-b@example.com")
	if err != nil || uB == nil {
		t.Fatal(err)
	}
	if uA.HouseholdID == uB.HouseholdID {
		t.Fatal("expected distinct households")
	}

	day := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)
	tid, err := st.CreateTransaction(ctx, uidA, uA.HouseholdID, -100, day, "mine", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = st.UpdateTransaction(ctx, uA.HouseholdID, uidB, tid, -200, day, "no", nil)
	if err == nil {
		t.Fatal("expected sql.ErrNoRows")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("got %v want %v", err, sql.ErrNoRows)
	}
}
