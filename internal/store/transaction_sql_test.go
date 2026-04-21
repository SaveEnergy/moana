package store

import "testing"

func TestSqlSelectOccurredAtAmountFromHouseholdInRange_matchesAppendOccurredAtShape(t *testing.T) {
	t.Parallel()
	want := sqlSelectOccurredAtAmountFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo
	if sqlSelectOccurredAtAmountFromHouseholdInRange != want {
		t.Fatalf("sqlSelectOccurredAtAmountFromHouseholdInRange must match prefix + from/to filters")
	}
}

func TestSqlSumIncomeExpenseInRange_matchesAppendOccurredAtShapes(t *testing.T) {
	t.Parallel()
	if g, w := sqlSumIncomeExpenseInRangeFromOnly, sqlSumIncomeExpenseBase+sqlFilterOccurredAtFrom; g != w {
		t.Fatalf("sqlSumIncomeExpenseInRangeFromOnly drift")
	}
	if g, w := sqlSumIncomeExpenseInRangeToOnly, sqlSumIncomeExpenseBase+sqlFilterOccurredAtTo; g != w {
		t.Fatalf("sqlSumIncomeExpenseInRangeToOnly drift")
	}
	if g, w := sqlSumIncomeExpenseInRangeBoth, sqlSumIncomeExpenseBase+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo; g != w {
		t.Fatalf("sqlSumIncomeExpenseInRangeBoth drift")
	}
}

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

func TestSqlTransactionInsertConditional_stable(t *testing.T) {
	t.Parallel()
	want := `
INSERT INTO transactions (user_id, amount_cents, occurred_at, description, category_id, created_at)
SELECT ?, ?, ?, ?, ?, ?
WHERE EXISTS (SELECT 1 FROM users WHERE id = ? AND household_id = ?)`
	if sqlTransactionInsertConditional != want {
		t.Fatalf("sqlTransactionInsertConditional drift")
	}
}

func TestSqlTransactionUpdateHouseholdScoped_stable(t *testing.T) {
	t.Parallel()
	want := `
UPDATE transactions SET amount_cents = ?, occurred_at = ?, description = ?, category_id = ?
WHERE id = ? AND user_id IN (SELECT id FROM users WHERE household_id = ?)
AND EXISTS (SELECT 1 FROM users WHERE id = ? AND household_id = ?)`
	if sqlTransactionUpdateHouseholdScoped != want {
		t.Fatalf("sqlTransactionUpdateHouseholdScoped drift")
	}
}
