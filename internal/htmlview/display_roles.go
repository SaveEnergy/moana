package htmlview

import "strings"

// AdminRoleLabel maps a role string to a display label.
func AdminRoleLabel(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin":
		return "Admin"
	default:
		return "Member"
	}
}

// HouseholdRoleLabel maps household membership role to a display label.
func HouseholdRoleLabel(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "owner":
		return "Owner"
	case "admin":
		return "Admin"
	default:
		return "Member"
	}
}
