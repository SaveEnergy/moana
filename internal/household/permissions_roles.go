package household

import (
	"strings"

	"moana/internal/store"
)

// CanManage reports whether the user may edit household metadata or members (owner/admin).
func CanManage(u *store.User) bool {
	return CanManageRole(u.HouseholdRole)
}

// CanManageRole is the role check without a full user struct.
func CanManageRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "owner", "admin":
		return true
	default:
		return false
	}
}

// CanLeave reports whether the user may leave the household (non-owners always; owner only when solo).
func CanLeave(u *store.User, memberCount int64) bool {
	if strings.ToLower(strings.TrimSpace(u.HouseholdRole)) != "owner" {
		return true
	}
	return memberCount == 1
}
