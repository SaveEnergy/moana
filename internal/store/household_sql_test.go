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

func TestSqlHouseholdInsertDefaultName_stable(t *testing.T) {
	t.Parallel()
	want := `INSERT INTO households (name, created_at) VALUES ('My household', ?)`
	if sqlHouseholdInsertDefaultName != want {
		t.Fatalf("sqlHouseholdInsertDefaultName drift")
	}
}

func TestSqlHouseholdUpdateName_stable(t *testing.T) {
	t.Parallel()
	want := `UPDATE households SET name = ? WHERE id = ?`
	if sqlHouseholdUpdateName != want {
		t.Fatalf("sqlHouseholdUpdateName drift")
	}
}
