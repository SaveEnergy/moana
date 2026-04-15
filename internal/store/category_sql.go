package store

// Full category row for list/get scans.
const sqlCategorySelectFull = `SELECT id, household_id, name, IFNULL(icon, ''), IFNULL(color, '') FROM categories`
