package store

// Full user row (credentials + household fields) for scanUser.
const sqlUserSelectFull = `SELECT id, email, password_hash, role, created_at, household_id, first_name, last_name, household_role FROM users`

// Insert user with empty names; household_role is typically 'owner' (new household) or 'member' (invite).
const sqlUserInsert = `INSERT INTO users (email, password_hash, role, created_at, household_id, first_name, last_name, household_role) VALUES (?, ?, ?, ?, ?, '', '', ?)`
