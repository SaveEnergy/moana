package store

// sqlHouseholdGetByID loads one household row by primary key.
const sqlHouseholdGetByID = `SELECT id, name, created_at FROM households WHERE id = ?`

// sqlHouseholdCountMembers counts users belonging to a household.
const sqlHouseholdCountMembers = `SELECT COUNT(*) FROM users WHERE household_id = ?`
