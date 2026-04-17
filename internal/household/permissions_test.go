package household

import (
	"testing"

	"moana/internal/store"
)

func TestCanManageRole(t *testing.T) {
	t.Parallel()
	if !CanManageRole("owner") || !CanManageRole("admin") || CanManageRole("member") || CanManageRole("") {
		t.Fatal()
	}
	if !CanManageRole("  ADMIN  ") || !CanManageRole("Owner") {
		t.Fatal("role check must be case-insensitive and trim whitespace")
	}
}

func TestCanManage(t *testing.T) {
	t.Parallel()
	if !CanManage(&store.User{HouseholdRole: "admin"}) {
		t.Fatal("admin should manage")
	}
	if CanManage(&store.User{HouseholdRole: "member"}) {
		t.Fatal("member should not manage")
	}
}

func TestCanRemoveMember(t *testing.T) {
	t.Parallel()
	hh := int64(42)
	owner := &store.User{HouseholdID: hh, HouseholdRole: "owner"}
	admin := &store.User{HouseholdID: hh, HouseholdRole: "admin"}
	member := &store.User{HouseholdID: hh, HouseholdRole: "member"}
	otherHH := &store.User{HouseholdID: 99, HouseholdRole: "member"}

	if !CanRemoveMember(owner, member) {
		t.Fatal("owner should remove member")
	}
	if CanRemoveMember(admin, owner) {
		t.Fatal("admin cannot remove owner")
	}
	if !CanRemoveMember(admin, member) {
		t.Fatal("admin should remove member")
	}
	if CanRemoveMember(owner, otherHH) {
		t.Fatal("different household")
	}
}

func TestCanLeave(t *testing.T) {
	t.Parallel()
	hh := int64(42)
	owner := &store.User{HouseholdID: hh, HouseholdRole: "owner"}
	member := &store.User{HouseholdID: hh, HouseholdRole: "member"}
	if !CanLeave(member, 99) {
		t.Fatal("non-owner may always leave")
	}
	if !CanLeave(owner, 1) {
		t.Fatal("sole owner may leave")
	}
	if CanLeave(owner, 2) {
		t.Fatal("owner with other members must transfer first")
	}
}

func TestCanLeave_ownerRoleTrimsWhitespace(t *testing.T) {
	t.Parallel()
	hh := int64(42)
	o := &store.User{HouseholdID: hh, HouseholdRole: "  owner  "}
	if CanLeave(o, 2) {
		t.Fatal("owner with peers cannot leave")
	}
	if !CanLeave(o, 1) {
		t.Fatal("sole owner can leave")
	}
}
