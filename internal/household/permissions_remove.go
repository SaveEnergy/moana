package household

import (
	"strings"

	"moana/internal/store"
)

// CanRemoveMember reports whether actor may remove target from the same household.
func CanRemoveMember(actor, target *store.User) bool {
	if actor.HouseholdID != target.HouseholdID {
		return false
	}
	if strings.ToLower(target.HouseholdRole) == "owner" {
		return false
	}
	ar := strings.ToLower(actor.HouseholdRole)
	tr := strings.ToLower(target.HouseholdRole)
	if ar == "owner" {
		return true
	}
	if ar == "admin" && tr == "member" {
		return true
	}
	return false
}

// CanRemoveMemberRow is for template helpers when only a [store.HouseholdMember] row is available.
func CanRemoveMemberRow(actor *store.User, member store.HouseholdMember) bool {
	target := &store.User{HouseholdID: actor.HouseholdID, HouseholdRole: member.HouseholdRole}
	return CanRemoveMember(actor, target)
}
