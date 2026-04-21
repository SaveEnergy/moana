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

func TestSqlSumAmountHouseholdRange_matchesAppendOccurredAtAndKindShapes(t *testing.T) {
	t.Parallel()
	if g, w := sqlSumAmountHouseholdRangeBoth, sqlSumAmountAllHousehold+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeBoth drift")
	}
	if g, w := sqlSumAmountHouseholdRangeFromOnly, sqlSumAmountAllHousehold+sqlFilterOccurredAtFrom; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeFromOnly drift")
	}
	if g, w := sqlSumAmountHouseholdRangeToOnly, sqlSumAmountAllHousehold+sqlFilterOccurredAtTo; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeToOnly drift")
	}
	if g, w := sqlSumAmountHouseholdRangeBothIncome, sqlSumAmountHouseholdRangeBoth+sqlFilterAmountIncome; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeBothIncome drift")
	}
	if g, w := sqlSumAmountHouseholdRangeBothExpense, sqlSumAmountHouseholdRangeBoth+sqlFilterAmountExpense; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeBothExpense drift")
	}
	if g, w := sqlSumAmountHouseholdRangeFromOnlyIncome, sqlSumAmountHouseholdRangeFromOnly+sqlFilterAmountIncome; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeFromOnlyIncome drift")
	}
	if g, w := sqlSumAmountHouseholdRangeFromOnlyExpense, sqlSumAmountHouseholdRangeFromOnly+sqlFilterAmountExpense; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeFromOnlyExpense drift")
	}
	if g, w := sqlSumAmountHouseholdRangeToOnlyIncome, sqlSumAmountHouseholdRangeToOnly+sqlFilterAmountIncome; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeToOnlyIncome drift")
	}
	if g, w := sqlSumAmountHouseholdRangeToOnlyExpense, sqlSumAmountHouseholdRangeToOnly+sqlFilterAmountExpense; g != w {
		t.Fatalf("sqlSumAmountHouseholdRangeToOnlyExpense drift")
	}
}

func TestSqlListTopExpenseCategoriesRange_matchesAppendOccurredAtShapes(t *testing.T) {
	t.Parallel()
	if g, w := sqlListTopExpenseCategoriesFromOnly, sqlListTopExpenseCategoriesPrefix+sqlFilterOccurredAtFrom+sqlListTopExpenseCategoriesSuffix; g != w {
		t.Fatalf("sqlListTopExpenseCategoriesFromOnly drift")
	}
	if g, w := sqlListTopExpenseCategoriesToOnly, sqlListTopExpenseCategoriesPrefix+sqlFilterOccurredAtTo+sqlListTopExpenseCategoriesSuffix; g != w {
		t.Fatalf("sqlListTopExpenseCategoriesToOnly drift")
	}
	if g, w := sqlListTopExpenseCategoriesBoth, sqlListTopExpenseCategoriesPrefix+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo+sqlListTopExpenseCategoriesSuffix; g != w {
		t.Fatalf("sqlListTopExpenseCategoriesBoth drift")
	}
}

func TestSqlTransactionListDated_matchesAppendOccurredAtAndKindShapes(t *testing.T) {
	t.Parallel()
	if g, w := sqlTransactionListDatedBothNoKindDescLimit, sqlTransactionListFromHousehold+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo+sqlTransactionListOrderDescLimit; g != w {
		t.Fatalf("sqlTransactionListDatedBothNoKindDescLimit drift")
	}
	if g, w := sqlTransactionListDatedFromOnlyIncomeAscLimit, sqlTransactionListFromHousehold+sqlFilterOccurredAtFrom+sqlFilterAmountIncome+sqlTransactionListOrderAscLimit; g != w {
		t.Fatalf("sqlTransactionListDatedFromOnlyIncomeAscLimit drift")
	}
	if g, w := sqlTransactionListDatedToOnlyExpenseDescLimit, sqlTransactionListFromHousehold+sqlFilterOccurredAtTo+sqlFilterAmountExpense+sqlTransactionListOrderDescLimit; g != w {
		t.Fatalf("sqlTransactionListDatedToOnlyExpenseDescLimit drift")
	}
	if g, w := sqlTransactionListDatedBothNoKindDesc, sqlTransactionListFromHousehold+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo+sqlTransactionListOrderDesc; g != w {
		t.Fatalf("sqlTransactionListDatedBothNoKindDesc drift")
	}
}

func TestSqlTransactionListSearchNoDate_matchesKindAndOrderShapes(t *testing.T) {
	t.Parallel()
	if g, w := sqlTransactionListSearchNoDateIncomeDescLimit, sqlTransactionListFromHousehold+sqlFilterAmountIncome+sqlTransactionListSearchLike+sqlTransactionListOrderDescLimit; g != w {
		t.Fatalf("sqlTransactionListSearchNoDateIncomeDescLimit drift")
	}
	if g, w := sqlTransactionListSearchNoDateNoKindDesc, sqlTransactionListFromHousehold+sqlTransactionListSearchLike+sqlTransactionListOrderDesc; g != w {
		t.Fatalf("sqlTransactionListSearchNoDateNoKindDesc drift")
	}
}

func TestSqlTransactionListDatedBothSearch_matchesFragmentOrder(t *testing.T) {
	t.Parallel()
	want := sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListSearchLike + sqlTransactionListOrderDescLimit
	if sqlTransactionListDatedBothSearchIncomeDescLimit != want {
		t.Fatalf("sqlTransactionListDatedBothSearchIncomeDescLimit drift")
	}
}

func TestSqlListCategoryAmountsRange_matchesAppendOccurredAtAndKindShapes(t *testing.T) {
	t.Parallel()
	if g, w := sqlListCategoryAmountsIncomeRangeBoth, sqlListCategoryAmountsSelectPrefix+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo+sqlFilterAmountIncome+sqlListCategoryAmountsGroupBy+sqlListCategoryAmountsOrderIncome; g != w {
		t.Fatalf("sqlListCategoryAmountsIncomeRangeBoth drift")
	}
	if g, w := sqlListCategoryAmountsIncomeRangeFromOnly, sqlListCategoryAmountsSelectPrefix+sqlFilterOccurredAtFrom+sqlFilterAmountIncome+sqlListCategoryAmountsGroupBy+sqlListCategoryAmountsOrderIncome; g != w {
		t.Fatalf("sqlListCategoryAmountsIncomeRangeFromOnly drift")
	}
	if g, w := sqlListCategoryAmountsIncomeRangeToOnly, sqlListCategoryAmountsSelectPrefix+sqlFilterOccurredAtTo+sqlFilterAmountIncome+sqlListCategoryAmountsGroupBy+sqlListCategoryAmountsOrderIncome; g != w {
		t.Fatalf("sqlListCategoryAmountsIncomeRangeToOnly drift")
	}
	if g, w := sqlListCategoryAmountsExpenseRangeBoth, sqlListCategoryAmountsSelectPrefix+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo+sqlFilterAmountExpense+sqlListCategoryAmountsGroupBy+sqlListCategoryAmountsOrderExpense; g != w {
		t.Fatalf("sqlListCategoryAmountsExpenseRangeBoth drift")
	}
	if g, w := sqlListCategoryAmountsExpenseRangeFromOnly, sqlListCategoryAmountsSelectPrefix+sqlFilterOccurredAtFrom+sqlFilterAmountExpense+sqlListCategoryAmountsGroupBy+sqlListCategoryAmountsOrderExpense; g != w {
		t.Fatalf("sqlListCategoryAmountsExpenseRangeFromOnly drift")
	}
	if g, w := sqlListCategoryAmountsExpenseRangeToOnly, sqlListCategoryAmountsSelectPrefix+sqlFilterOccurredAtTo+sqlFilterAmountExpense+sqlListCategoryAmountsGroupBy+sqlListCategoryAmountsOrderExpense; g != w {
		t.Fatalf("sqlListCategoryAmountsExpenseRangeToOnly drift")
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
