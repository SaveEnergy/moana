package household

import (
	"context"
	"testing"

	"moana/internal/dbutil"
	"moana/internal/passwordtest"
)

func TestLoadSettingsPage_memberCountMatchesMemberList(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
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

	hash2 := passwordtest.MustHash(t, "y")
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
