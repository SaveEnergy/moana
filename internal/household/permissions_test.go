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
