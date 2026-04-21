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

func TestSqlCategoryInsert_stable(t *testing.T) {
	t.Parallel()
	want := `INSERT INTO categories (household_id, name, icon, color) VALUES (?, ?, ?, ?)`
	if sqlCategoryInsert != want {
		t.Fatalf("sqlCategoryInsert drift")
	}
}

func TestSqlCategoryUpdate_stable(t *testing.T) {
	t.Parallel()
	want := `UPDATE categories SET name = ?, icon = ?, color = ? WHERE id = ? AND household_id = ?`
	if sqlCategoryUpdate != want {
		t.Fatalf("sqlCategoryUpdate drift")
	}
}

func TestSqlCategoryDelete_stable(t *testing.T) {
	t.Parallel()
	want := `DELETE FROM categories WHERE id = ? AND household_id = ?`
	if sqlCategoryDelete != want {
		t.Fatalf("sqlCategoryDelete drift")
	}
}
