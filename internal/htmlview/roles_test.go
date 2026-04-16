package htmlview

import "testing"

func TestAdminRoleLabel(t *testing.T) {
	t.Parallel()
	if got := AdminRoleLabel("  ADMIN  "); got != "Admin" {
		t.Fatalf("got %q", got)
	}
	if got := AdminRoleLabel("user"); got != "Member" {
		t.Fatalf("got %q", got)
	}
}

func TestHouseholdRoleLabel(t *testing.T) {
	t.Parallel()
	if got := HouseholdRoleLabel("owner"); got != "Owner" {
		t.Fatalf("got %q", got)
	}
	if got := HouseholdRoleLabel("admin"); got != "Admin" {
		t.Fatalf("got %q", got)
	}
	if got := HouseholdRoleLabel("member"); got != "Member" {
		t.Fatalf("got %q", got)
	}
}
