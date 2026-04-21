package store

// Full category row for list/get scans.
const sqlCategorySelectFull = `SELECT id, household_id, name, IFNULL(icon, ''), IFNULL(color, '') FROM categories`

// sqlCategoryListByHousehold lists categories for one household (case-insensitive name order).
const sqlCategoryListByHousehold = sqlCategorySelectFull + ` WHERE household_id = ? ORDER BY name COLLATE NOCASE`

// sqlCategoryGetByIDHousehold loads one category row scoped to household.
const sqlCategoryGetByIDHousehold = sqlCategorySelectFull + ` WHERE id = ? AND household_id = ?`
