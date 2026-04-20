package store

import (
	"context"
	"testing"

	"moana/internal/passwordtest"
)

func TestGetCategoryByID_found(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "gcid-found@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	catID, err := st.CreateCategory(ctx, u.HouseholdID, "Food", "", "#112233")
	if err != nil {
		t.Fatal(err)
	}
	got, err := st.GetCategoryByID(ctx, u.HouseholdID, catID)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != catID || got.Name != "Food" || got.HouseholdID != u.HouseholdID {
		t.Fatalf("category %+v", got)
	}
}

func TestGetCategoryByID_notFound(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "gcid-nf@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	got, err := st.GetCategoryByID(ctx, u.HouseholdID, 999999999)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("want nil, got %+v", got)
	}
}

func TestGetCategoryByID_wrongHouseholdReturnsNil(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uidA, err := st.CreateUser(ctx, "gcid-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uA, err := st.GetUserByID(ctx, uidA)
	if err != nil || uA == nil {
		t.Fatal(err)
	}
	catID, err := st.CreateCategory(ctx, uA.HouseholdID, "Other HH", "", "")
	if err != nil {
		t.Fatal(err)
	}

	uidB, err := st.CreateUser(ctx, "gcid-b@example.com", hash, "user")
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

	got, err := st.GetCategoryByID(ctx, uB.HouseholdID, catID)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("category from another household must not resolve: %+v", got)
	}
}

func TestListCategories_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListCategories(alreadyCancelledContext(t), 1)
	assertErrIsContextCanceled(t, err)
}

func TestGetCategoryByID_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.GetCategoryByID(alreadyCancelledContext(t), 1, 1)
	assertErrIsContextCanceled(t, err)
}
