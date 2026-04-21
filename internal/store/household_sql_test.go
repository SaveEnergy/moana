package store

import "testing"

func TestSqlHouseholdGetByID_stable(t *testing.T) {
	t.Parallel()
	want := `SELECT id, name, created_at FROM households WHERE id = ?`
	if sqlHouseholdGetByID != want {
		t.Fatalf("sqlHouseholdGetByID drift")
	}
}

func TestSqlHouseholdCountMembers_stable(t *testing.T) {
	t.Parallel()
	want := `SELECT COUNT(*) FROM users WHERE household_id = ?`
	if sqlHouseholdCountMembers != want {
		t.Fatalf("sqlHouseholdCountMembers drift")
	}
}
