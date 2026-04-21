package store

// Household-scoped transaction rows joined to category + owner user (for visibility).
const sqlTransactionSelectFromHousehold = `
SELECT t.id, t.user_id, t.amount_cents, t.occurred_at, t.description, t.category_id, COALESCE(c.name, ''), COALESCE(IFNULL(c.icon, ''), ''), t.created_at
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
INNER JOIN users owner ON owner.id = t.user_id`

// sqlTransactionListFromHousehold is [sqlTransactionSelectFromHousehold] plus the household scope
// WHERE (filters are appended: date bounds, kind, search, order, limit).
const sqlTransactionListFromHousehold = sqlTransactionSelectFromHousehold + `
WHERE owner.household_id = ?`

// sqlTransactionListOrderDesc / Asc are shared by [ListTransactions] static queries and the dynamic SQL builder.
const sqlTransactionListOrderDesc = ` ORDER BY t.occurred_at DESC, t.id DESC`
const sqlTransactionListOrderAsc = ` ORDER BY t.occurred_at ASC, t.id ASC`

const sqlTransactionListLimitSuffix = ` LIMIT ?`

const sqlTransactionListOrderDescLimit = sqlTransactionListOrderDesc + sqlTransactionListLimitSuffix
const sqlTransactionListOrderAscLimit = sqlTransactionListOrderAsc + sqlTransactionListLimitSuffix

// sqlTransactionListRecentDescLimit is [ListTransactions] when only Limit is set (default sort, no filters).
const sqlTransactionListRecentDescLimit = sqlTransactionListFromHousehold + sqlTransactionListOrderDescLimit

// sqlTransactionListRecentDescLimitIncome / Expense are [ListTransactions] with kind income|expense, Limit, and no other filters.
const sqlTransactionListRecentDescLimitIncome = sqlTransactionListFromHousehold + sqlFilterAmountIncome + sqlTransactionListOrderDescLimit

const sqlTransactionListRecentDescLimitExpense = sqlTransactionListFromHousehold + sqlFilterAmountExpense + sqlTransactionListOrderDescLimit

// sqlTransactionListOldestAscLimit variants are [ListTransactions] with OldestFirst, Limit, and no date/search filters.
const sqlTransactionListOldestAscLimit = sqlTransactionListFromHousehold + sqlTransactionListOrderAscLimit

const sqlTransactionListOldestAscLimitIncome = sqlTransactionListFromHousehold + sqlFilterAmountIncome + sqlTransactionListOrderAscLimit

const sqlTransactionListOldestAscLimitExpense = sqlTransactionListFromHousehold + sqlFilterAmountExpense + sqlTransactionListOrderAscLimit

// sqlTransactionListDated* are [ListTransactions] with optional occurred_at bounds, no search, LIMIT ? (same fragment order as the dynamic builder).
const (
	sqlTransactionListDatedBothNoKindDescLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedBothNoKindAscLimit   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlTransactionListOrderAscLimit
	sqlTransactionListDatedBothIncomeDescLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedBothIncomeAscLimit   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderAscLimit
	sqlTransactionListDatedBothExpenseDescLimit = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedBothExpenseAscLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderAscLimit

	sqlTransactionListDatedFromOnlyNoKindDescLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedFromOnlyNoKindAscLimit   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlTransactionListOrderAscLimit
	sqlTransactionListDatedFromOnlyIncomeDescLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountIncome + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedFromOnlyIncomeAscLimit   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountIncome + sqlTransactionListOrderAscLimit
	sqlTransactionListDatedFromOnlyExpenseDescLimit = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountExpense + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedFromOnlyExpenseAscLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountExpense + sqlTransactionListOrderAscLimit

	sqlTransactionListDatedToOnlyNoKindDescLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedToOnlyNoKindAscLimit   = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlTransactionListOrderAscLimit
	sqlTransactionListDatedToOnlyIncomeDescLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedToOnlyIncomeAscLimit   = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderAscLimit
	sqlTransactionListDatedToOnlyExpenseDescLimit = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderDescLimit
	sqlTransactionListDatedToOnlyExpenseAscLimit  = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderAscLimit
)

// sqlTransactionListDated*NoLimit are [ListTransactions] with date bounds and no search, without LIMIT (unbounded scan for the filter).
const (
	sqlTransactionListDatedBothNoKindDesc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlTransactionListOrderDesc
	sqlTransactionListDatedBothNoKindAsc   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlTransactionListOrderAsc
	sqlTransactionListDatedBothIncomeDesc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderDesc
	sqlTransactionListDatedBothIncomeAsc   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderAsc
	sqlTransactionListDatedBothExpenseDesc = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderDesc
	sqlTransactionListDatedBothExpenseAsc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderAsc

	sqlTransactionListDatedFromOnlyNoKindDesc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlTransactionListOrderDesc
	sqlTransactionListDatedFromOnlyNoKindAsc   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlTransactionListOrderAsc
	sqlTransactionListDatedFromOnlyIncomeDesc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountIncome + sqlTransactionListOrderDesc
	sqlTransactionListDatedFromOnlyIncomeAsc   = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountIncome + sqlTransactionListOrderAsc
	sqlTransactionListDatedFromOnlyExpenseDesc = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountExpense + sqlTransactionListOrderDesc
	sqlTransactionListDatedFromOnlyExpenseAsc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterAmountExpense + sqlTransactionListOrderAsc

	sqlTransactionListDatedToOnlyNoKindDesc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlTransactionListOrderDesc
	sqlTransactionListDatedToOnlyNoKindAsc   = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlTransactionListOrderAsc
	sqlTransactionListDatedToOnlyIncomeDesc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderDesc
	sqlTransactionListDatedToOnlyIncomeAsc   = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlTransactionListOrderAsc
	sqlTransactionListDatedToOnlyExpenseDesc = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderDesc
	sqlTransactionListDatedToOnlyExpenseAsc  = sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderAsc
)

// sqlTransactionListSearchLike is appended after household/date/kind filters for [ListTransactions] search (two LIKE placeholders).
const sqlTransactionListSearchLike = ` AND (t.description LIKE ? ESCAPE '!' OR COALESCE(c.name, '') LIKE ? ESCAPE '!')`

// sqlTransactionGetByIDHousehold loads one row by transaction id, scoped to the owner household.
const sqlTransactionGetByIDHousehold = sqlTransactionSelectFromHousehold + `
WHERE t.id = ? AND owner.household_id = ?`

// sqlFromHouseholdTx scopes transactions to a household without joining category columns (sums, movement).
const sqlFromHouseholdTx = `
FROM transactions t
INNER JOIN users owner ON owner.id = t.user_id
WHERE owner.household_id = ?`

// sqlSelectOccurredAtAmountFromHousehold is the SELECT + FROM/WHERE prefix for per-row occurred_at + amount reads.
const sqlSelectOccurredAtAmountFromHousehold = `SELECT t.occurred_at, t.amount_cents ` + sqlFromHouseholdTx

// sqlSelectOccurredAtAmountFromHouseholdInRange is [DailyAbsMovementByLocalDate]: household id plus inclusive UTC bounds on occurred_at (same shape as [appendOccurredAtRange] with both bounds set).
const sqlSelectOccurredAtAmountFromHouseholdInRange = sqlSelectOccurredAtAmountFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo

// sqlSumAmountAllHousehold is [SumAmountCentsByKind] with no date bounds and no income/expense sign filter (kind "").
const sqlSumAmountAllHousehold = `SELECT COALESCE(SUM(t.amount_cents), 0) ` + sqlFromHouseholdTx

// sqlSumAmountIncomeHouseholdNoBounds / sqlSumAmountExpenseHouseholdNoBounds are [SumAmountCentsByKind] with no date bounds and kind income|expense.
const (
	sqlSumAmountIncomeHouseholdNoBounds  = sqlSumAmountAllHousehold + sqlFilterAmountIncome
	sqlSumAmountExpenseHouseholdNoBounds = sqlSumAmountAllHousehold + sqlFilterAmountExpense
)

// sqlSumAmountHouseholdRange* are [SumAmountCentsByKind] with date bounds (income|expense filters after occurred_at, same order as [appendOccurredAtRange] + [sqlAmountKindFilter]).
const (
	sqlSumAmountHouseholdRangeBoth     = sqlSumAmountAllHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo
	sqlSumAmountHouseholdRangeFromOnly = sqlSumAmountAllHousehold + sqlFilterOccurredAtFrom
	sqlSumAmountHouseholdRangeToOnly   = sqlSumAmountAllHousehold + sqlFilterOccurredAtTo

	sqlSumAmountHouseholdRangeBothIncome      = sqlSumAmountHouseholdRangeBoth + sqlFilterAmountIncome
	sqlSumAmountHouseholdRangeBothExpense     = sqlSumAmountHouseholdRangeBoth + sqlFilterAmountExpense
	sqlSumAmountHouseholdRangeFromOnlyIncome  = sqlSumAmountHouseholdRangeFromOnly + sqlFilterAmountIncome
	sqlSumAmountHouseholdRangeFromOnlyExpense = sqlSumAmountHouseholdRangeFromOnly + sqlFilterAmountExpense
	sqlSumAmountHouseholdRangeToOnlyIncome    = sqlSumAmountHouseholdRangeToOnly + sqlFilterAmountIncome
	sqlSumAmountHouseholdRangeToOnlyExpense   = sqlSumAmountHouseholdRangeToOnly + sqlFilterAmountExpense
)

// sqlSumIncomeExpenseBase is [SumIncomeExpenseCentsInRange] with no date bounds (appendOccurredAtRange adds optional ANDs).
const sqlSumIncomeExpenseBase = `SELECT COALESCE(SUM(CASE WHEN t.amount_cents > 0 THEN t.amount_cents ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN t.amount_cents < 0 THEN t.amount_cents ELSE 0 END), 0) ` + sqlFromHouseholdTx

// sqlSumIncomeExpenseInRange* are [SumIncomeExpenseCentsInRange] with optional occurred_at bounds (same fragments as [appendOccurredAtRange]).
const (
	sqlSumIncomeExpenseInRangeFromOnly = sqlSumIncomeExpenseBase + sqlFilterOccurredAtFrom
	sqlSumIncomeExpenseInRangeToOnly   = sqlSumIncomeExpenseBase + sqlFilterOccurredAtTo
	sqlSumIncomeExpenseInRangeBoth     = sqlSumIncomeExpenseBase + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo
)

// sqlSumRunningTotalAndIncomeExpenseInTwoRanges is the dashboard aggregate: one scan for all-time net plus
// income/expense in two closed ranges. Placeholders: aFrom, aTo (each range uses two CASE branches); then household_id.
const sqlSumRunningTotalAndIncomeExpenseInTwoRanges = `SELECT
  COALESCE(SUM(t.amount_cents), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents > 0 THEN t.amount_cents ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents < 0 THEN t.amount_cents ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents > 0 THEN t.amount_cents ELSE 0 END), 0),
  COALESCE(SUM(CASE WHEN t.occurred_at >= ? AND t.occurred_at <= ? AND t.amount_cents < 0 THEN t.amount_cents ELSE 0 END), 0)
` + sqlFromHouseholdTx

// sqlAggregateFromHouseholdTx is FROM/WHERE for household aggregates that need the category join (matches listing joins).
const sqlAggregateFromHouseholdTx = `
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
INNER JOIN users owner ON owner.id = t.user_id
WHERE owner.household_id = ?`

// sqlListTopExpenseCategoriesPrefix and suffix wrap optional [appendOccurredAtRange] fragments for [ListTopExpenseCategories].
const sqlListTopExpenseCategoriesPrefix = `SELECT t.category_id, COALESCE(c.name, 'Uncategorized'), COALESCE(SUM(t.amount_cents), 0)
` + sqlAggregateFromHouseholdTx + sqlFilterAmountExpense

const sqlListTopExpenseCategoriesSuffix = ` GROUP BY t.category_id, c.name ORDER BY SUM(t.amount_cents) ASC LIMIT ?`

// sqlListTopExpenseCategoriesNoDate is [ListTopExpenseCategories] with no date bounds.
const sqlListTopExpenseCategoriesNoDate = sqlListTopExpenseCategoriesPrefix + sqlListTopExpenseCategoriesSuffix

// sqlListTopExpenseCategoriesFromOnly / ToOnly / Both are [ListTopExpenseCategories] with optional occurred_at bounds (same fragments as [appendOccurredAtRange]).
const (
	sqlListTopExpenseCategoriesFromOnly = sqlListTopExpenseCategoriesPrefix + sqlFilterOccurredAtFrom + sqlListTopExpenseCategoriesSuffix
	sqlListTopExpenseCategoriesToOnly   = sqlListTopExpenseCategoriesPrefix + sqlFilterOccurredAtTo + sqlListTopExpenseCategoriesSuffix
	sqlListTopExpenseCategoriesBoth     = sqlListTopExpenseCategoriesPrefix + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlListTopExpenseCategoriesSuffix
)

// sqlListCategoryAmountsSelectPrefix is the SELECT + FROM for [ListCategoryAmountsInRange] (date + kind filters appended).
const sqlListCategoryAmountsSelectPrefix = `SELECT t.category_id, COALESCE(MAX(c.name), 'Uncategorized'), COALESCE(MAX(IFNULL(c.icon, '')), ''), COALESCE(MAX(IFNULL(c.color, '')), ''), COALESCE(SUM(t.amount_cents), 0)
` + sqlAggregateFromHouseholdTx

const sqlListCategoryAmountsGroupBy = ` GROUP BY t.category_id`

const sqlListCategoryAmountsOrderIncome = ` ORDER BY SUM(t.amount_cents) DESC`
const sqlListCategoryAmountsOrderExpense = ` ORDER BY SUM(t.amount_cents) ASC`

// sqlListCategoryAmountsIncomeFullHousehold / ExpenseFullHousehold are [ListCategoryAmountsInRange] with no date bounds.
const sqlListCategoryAmountsIncomeFullHousehold = sqlListCategoryAmountsSelectPrefix + sqlFilterAmountIncome + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderIncome

const sqlListCategoryAmountsExpenseFullHousehold = sqlListCategoryAmountsSelectPrefix + sqlFilterAmountExpense + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderExpense

// sqlListCategoryAmountsIncomeRange* / ExpenseRange* are [ListCategoryAmountsInRange] with date bounds and kind (same fragment order as dynamic builder).
const (
	sqlListCategoryAmountsIncomeRangeBoth     = sqlListCategoryAmountsSelectPrefix + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderIncome
	sqlListCategoryAmountsIncomeRangeFromOnly = sqlListCategoryAmountsSelectPrefix + sqlFilterOccurredAtFrom + sqlFilterAmountIncome + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderIncome
	sqlListCategoryAmountsIncomeRangeToOnly   = sqlListCategoryAmountsSelectPrefix + sqlFilterOccurredAtTo + sqlFilterAmountIncome + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderIncome

	sqlListCategoryAmountsExpenseRangeBoth     = sqlListCategoryAmountsSelectPrefix + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderExpense
	sqlListCategoryAmountsExpenseRangeFromOnly = sqlListCategoryAmountsSelectPrefix + sqlFilterOccurredAtFrom + sqlFilterAmountExpense + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderExpense
	sqlListCategoryAmountsExpenseRangeToOnly   = sqlListCategoryAmountsSelectPrefix + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderExpense
)

// sqlTransactionInsertConditional inserts a row only if user_id belongs to household_id.
const sqlTransactionInsertConditional = `
INSERT INTO transactions (user_id, amount_cents, occurred_at, description, category_id, created_at)
SELECT ?, ?, ?, ?, ?, ?
WHERE EXISTS (SELECT 1 FROM users WHERE id = ? AND household_id = ?)`

// sqlTransactionUpdateHouseholdScoped updates a row scoped to household membership ([UpdateTransaction]).
const sqlTransactionUpdateHouseholdScoped = `
UPDATE transactions SET amount_cents = ?, occurred_at = ?, description = ?, category_id = ?
WHERE id = ? AND user_id IN (SELECT id FROM users WHERE household_id = ?)
AND EXISTS (SELECT 1 FROM users WHERE id = ? AND household_id = ?)`
