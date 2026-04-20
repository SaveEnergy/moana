package store

import (
	"context"
	"testing"

	"moana/internal/passwordtest"
)

func TestUpdateHouseholdName(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "uhname@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	if err := st.UpdateHouseholdName(ctx, hid, "Renamed Household"); err != nil {
		t.Fatal(err)
	}
	h, err := st.GetHousehold(ctx, hid)
	if err != nil {
		t.Fatal(err)
	}
	if h == nil || h.Name != "Renamed Household" {
		t.Fatalf("household %+v", h)
	}
}
