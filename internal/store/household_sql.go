package store

// sqlHouseholdGetByID loads one household row by primary key.
const sqlHouseholdGetByID = `SELECT id, name, created_at FROM households WHERE id = ?`

// sqlHouseholdCountMembers counts users belonging to a household.
const sqlHouseholdCountMembers = `SELECT COUNT(*) FROM users WHERE household_id = ?`

// sqlHouseholdInsertDefaultName inserts a new household with the default display name.
const sqlHouseholdInsertDefaultName = `INSERT INTO households (name, created_at) VALUES ('My household', ?)`

// sqlHouseholdUpdateName sets the household display name.
const sqlHouseholdUpdateName = `UPDATE households SET name = ? WHERE id = ?`
