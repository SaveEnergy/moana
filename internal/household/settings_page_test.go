package household

import (
	"context"
	"testing"

	"moana/internal/auth"
	"moana/internal/dbutil"
)

func TestLoadSettingsPage_memberCountMatchesMemberList(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "settings-owner@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, uid)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hid := owner.HouseholdID

	data, err := LoadSettingsPage(ctx, st, owner, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if data.MemberCount != 1 || len(data.Members) != 1 {
		t.Fatalf("solo household: MemberCount=%d len=%d", data.MemberCount, len(data.Members))
	}

	hash2, err := auth.HashPassword("y")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateHouseholdMember(ctx, hid, "settings-member@example.com", hash2); err != nil {
		t.Fatal(err)
	}
	data2, err := LoadSettingsPage(ctx, st, owner, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if data2.MemberCount != 2 || len(data2.Members) != 2 {
		t.Fatalf("two members: MemberCount=%d len=%d", data2.MemberCount, len(data2.Members))
	}
}
