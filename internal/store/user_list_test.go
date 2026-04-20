package store

import (
	"context"
	"testing"

	"moana/internal/passwordtest"
)

func TestListHouseholdMembers_orderAndRoles(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	ownerID, err := st.CreateUser(ctx, "lhs-owner@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, ownerID)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hid := owner.HouseholdID

	hash2 := passwordtest.MustHash(t, "y")
	_, err = st.CreateHouseholdMember(ctx, hid, "lhs-member@example.com", hash2)
	if err != nil {
		t.Fatal(err)
	}

	members, err := st.ListHouseholdMembers(ctx, hid)
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 2 {
		t.Fatalf("len %d want 2", len(members))
	}
	if members[0].ID != ownerID {
		t.Fatalf("first row id %d want owner id %d", members[0].ID, ownerID)
	}
	if members[0].Email != "lhs-owner@example.com" || members[1].Email != "lhs-member@example.com" {
		t.Fatalf("emails %+v %+v", members[0], members[1])
	}
	if members[0].HouseholdRole != "owner" || members[1].HouseholdRole != "member" {
		t.Fatalf("roles %+v %+v", members[0].HouseholdRole, members[1].HouseholdRole)
	}
}
