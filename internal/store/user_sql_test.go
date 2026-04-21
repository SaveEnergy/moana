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
