package store

// Full user row (credentials + household fields) for scanUser.
const sqlUserSelectFull = `SELECT id, email, password_hash, role, created_at, household_id, first_name, last_name, household_role FROM users`

// sqlUserGetByEmailCaseInsensitive loads one user by email with COLLATE NOCASE (see [normalizeUserEmail]).
const sqlUserGetByEmailCaseInsensitive = sqlUserSelectFull + ` WHERE email = ? COLLATE NOCASE`

// sqlUserGetByID loads one user by primary key.
const sqlUserGetByID = sqlUserSelectFull + ` WHERE id = ?`

// sqlUserListSummaryAdmin lists accounts for admin views (subset of columns).
const sqlUserListSummaryAdmin = `SELECT id, email, role, created_at FROM users ORDER BY id`

// sqlUserListHouseholdMembers lists members of one household (display fields + role).
const sqlUserListHouseholdMembers = `SELECT id, email, IFNULL(first_name, ''), IFNULL(last_name, ''), household_role
FROM users WHERE household_id = ? ORDER BY id`

// sqlUserUpdatePassword sets the password hash for one user.
const sqlUserUpdatePassword = `UPDATE users SET password_hash = ? WHERE id = ?`

// sqlUserUpdateProfile sets first/last name for one user.
const sqlUserUpdateProfile = `UPDATE users SET first_name = ?, last_name = ? WHERE id = ?`

// sqlUserDetachToSoloHousehold moves a user to a new household as owner ([DetachUserToSoloHousehold]).
const sqlUserDetachToSoloHousehold = `UPDATE users SET household_id = ?, household_role = 'owner' WHERE id = ?`

// Insert user with empty names; household_role is typically 'owner' (new household) or 'member' (invite).
const sqlUserInsert = `INSERT INTO users (email, password_hash, role, created_at, household_id, first_name, last_name, household_role) VALUES (?, ?, ?, ?, ?, '', '', ?)`
