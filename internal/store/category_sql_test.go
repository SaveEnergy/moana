package store

import "testing"

func TestSqlCategoryListByHousehold_matchesConcatenation(t *testing.T) {
	t.Parallel()
	want := sqlCategorySelectFull + " WHERE household_id = ? ORDER BY name COLLATE NOCASE"
	if sqlCategoryListByHousehold != want {
		t.Fatalf("sqlCategoryListByHousehold must stay aligned with select + household list ORDER BY")
	}
}

func TestSqlCategoryGetByIDHousehold_matchesConcatenation(t *testing.T) {
	t.Parallel()
	want := sqlCategorySelectFull + " WHERE id = ? AND household_id = ?"
	if sqlCategoryGetByIDHousehold != want {
		t.Fatalf("sqlCategoryGetByIDHousehold must stay aligned with select + id/household WHERE")
	}
}
