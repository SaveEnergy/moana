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

// sqlTransactionGetByIDHousehold loads one row by transaction id, scoped to the owner household.
const sqlTransactionGetByIDHousehold = sqlTransactionSelectFromHousehold + `
WHERE t.id = ? AND owner.household_id = ?`

// sqlFromHouseholdTx scopes transactions to a household without joining category columns (sums, movement).
const sqlFromHouseholdTx = `
FROM transactions t
INNER JOIN users owner ON owner.id = t.user_id
WHERE owner.household_id = ?`

// sqlSumIncomeExpenseBase is [SumIncomeExpenseCentsInRange] with no date bounds (appendOccurredAtRange adds optional ANDs).
const sqlSumIncomeExpenseBase = `SELECT COALESCE(SUM(CASE WHEN t.amount_cents > 0 THEN t.amount_cents ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN t.amount_cents < 0 THEN t.amount_cents ELSE 0 END), 0) ` + sqlFromHouseholdTx

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
