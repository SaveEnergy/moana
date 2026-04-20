package store

import (
	"context"
	"testing"

	"moana/internal/passwordtest"
)

func TestGetHousehold_found(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "gh-found@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	h, err := st.GetHousehold(ctx, u.HouseholdID)
	if err != nil {
		t.Fatal(err)
	}
	if h == nil || h.ID != u.HouseholdID {
		t.Fatalf("household %+v", h)
	}
	if h.Name != "My household" {
		t.Fatalf("name %q", h.Name)
	}
	if h.CreatedAt.IsZero() {
		t.Fatal("CreatedAt unset")
	}
}

func TestGetHousehold_notFound(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	h, err := st.GetHousehold(ctx, 999999)
	if err != nil {
		t.Fatal(err)
	}
	if h != nil {
		t.Fatalf("want nil, got %+v", h)
	}
}

func TestCountHouseholdMembers(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "cnt-owner@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, uid)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hid := owner.HouseholdID

	n, err := st.CountHouseholdMembers(ctx, hid)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("count %d want 1", n)
	}

	hash2 := passwordtest.MustHash(t, "y")
	if _, err := st.CreateHouseholdMember(ctx, hid, "cnt-mem@example.com", hash2); err != nil {
		t.Fatal(err)
	}
	n, err = st.CountHouseholdMembers(ctx, hid)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("count %d want 2", n)
	}
}

func TestGetHousehold_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.GetHousehold(alreadyCancelledContext(t), 1)
	assertErrIsContextCanceled(t, err)
}

func TestCountHouseholdMembers_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.CountHouseholdMembers(alreadyCancelledContext(t), 1)
	assertErrIsContextCanceled(t, err)
}
