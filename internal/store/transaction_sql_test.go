package store

import "testing"

func TestSqlTransactionListFromHousehold_matchesSelectPlusHouseholdWhere(t *testing.T) {
	t.Parallel()
	want := sqlTransactionSelectFromHousehold + "\nWHERE owner.household_id = ?"
	if sqlTransactionListFromHousehold != want {
		t.Fatalf("sqlTransactionListFromHousehold must stay aligned with select + household WHERE placeholder")
	}
}

func TestSqlTransactionGetByIDHousehold_matchesSelectPlusIdAndHouseholdWhere(t *testing.T) {
	t.Parallel()
	want := sqlTransactionSelectFromHousehold + "\nWHERE t.id = ? AND owner.household_id = ?"
	if sqlTransactionGetByIDHousehold != want {
		t.Fatalf("sqlTransactionGetByIDHousehold must stay aligned with select + id/household WHERE placeholders")
	}
}
