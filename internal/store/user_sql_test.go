package store

import "testing"

func TestSqlUserGetByEmailCaseInsensitive_matchesConcatenation(t *testing.T) {
	t.Parallel()
	want := sqlUserSelectFull + " WHERE email = ? COLLATE NOCASE"
	if sqlUserGetByEmailCaseInsensitive != want {
		t.Fatalf("sqlUserGetByEmailCaseInsensitive must stay aligned with select + email WHERE")
	}
}

func TestSqlUserGetByID_matchesConcatenation(t *testing.T) {
	t.Parallel()
	want := sqlUserSelectFull + " WHERE id = ?"
	if sqlUserGetByID != want {
		t.Fatalf("sqlUserGetByID must stay aligned with select + id WHERE")
	}
}

func TestSqlUserListSummaryAdmin_stable(t *testing.T) {
	t.Parallel()
	want := `SELECT id, email, role, created_at FROM users ORDER BY id`
	if sqlUserListSummaryAdmin != want {
		t.Fatalf("sqlUserListSummaryAdmin drift")
	}
}

func TestSqlUserListHouseholdMembers_stable(t *testing.T) {
	t.Parallel()
	want := `SELECT id, email, IFNULL(first_name, ''), IFNULL(last_name, ''), household_role
FROM users WHERE household_id = ? ORDER BY id`
	if sqlUserListHouseholdMembers != want {
		t.Fatalf("sqlUserListHouseholdMembers drift")
	}
}

func TestSqlUserUpdatePassword_stable(t *testing.T) {
	t.Parallel()
	want := `UPDATE users SET password_hash = ? WHERE id = ?`
	if sqlUserUpdatePassword != want {
		t.Fatalf("sqlUserUpdatePassword drift")
	}
}

func TestSqlUserUpdateProfile_stable(t *testing.T) {
	t.Parallel()
	want := `UPDATE users SET first_name = ?, last_name = ? WHERE id = ?`
	if sqlUserUpdateProfile != want {
		t.Fatalf("sqlUserUpdateProfile drift")
	}
}

func TestSqlUserDetachToSoloHousehold_stable(t *testing.T) {
	t.Parallel()
	want := `UPDATE users SET household_id = ?, household_role = 'owner' WHERE id = ?`
	if sqlUserDetachToSoloHousehold != want {
		t.Fatalf("sqlUserDetachToSoloHousehold drift")
	}
}
