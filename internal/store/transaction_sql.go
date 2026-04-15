package store

// Household-scoped transaction rows joined to category + owner user (for visibility).
const sqlTransactionSelectFromHousehold = `
SELECT t.id, t.user_id, t.amount_cents, t.occurred_at, t.description, t.category_id, COALESCE(c.name, ''), COALESCE(IFNULL(c.icon, ''), ''), t.created_at
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
INNER JOIN users owner ON owner.id = t.user_id`

// sqlFromHouseholdTx scopes transactions to a household without joining category columns (sums, movement).
const sqlFromHouseholdTx = `
FROM transactions t
INNER JOIN users owner ON owner.id = t.user_id
WHERE owner.household_id = ?`

// sqlAggregateFromHouseholdTx is FROM/WHERE for household aggregates that need the category join (matches listing joins).
const sqlAggregateFromHouseholdTx = `
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
INNER JOIN users owner ON owner.id = t.user_id
WHERE owner.household_id = ?`
