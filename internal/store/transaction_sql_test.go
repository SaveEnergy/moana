package store

import "testing"

func TestSqlTransactionListFromHousehold_matchesSelectPlusHouseholdWhere(t *testing.T) {
	t.Parallel()
	want := sqlTransactionSelectFromHousehold + "\nWHERE owner.household_id = ?"
	if sqlTransactionListFromHousehold != want {
		t.Fatalf("sqlTransactionListFromHousehold must stay aligned with select + household WHERE placeholder")
	}
}
