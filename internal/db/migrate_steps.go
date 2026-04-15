package db

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"moana/internal/timeutil"
)

func migrateV1(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE COLLATE NOCASE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('user', 'admin')),
  created_at TEXT NOT NULL
);
CREATE TABLE categories (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  UNIQUE (user_id, name)
);
CREATE TABLE transactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount_cents INTEGER NOT NULL,
  occurred_at TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
  created_at TEXT NOT NULL
);
CREATE INDEX idx_transactions_user_occurred ON transactions(user_id, occurred_at);
CREATE INDEX idx_categories_user ON categories(user_id);
`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (1)`)
	return err
}

func migrateV2(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `ALTER TABLE categories ADD COLUMN icon TEXT NOT NULL DEFAULT '';`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (2)`)
	return err
}

func migrateV3(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `
CREATE TABLE households (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL DEFAULT 'My household',
  created_at TEXT NOT NULL
);`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `ALTER TABLE users ADD COLUMN household_id INTEGER REFERENCES households(id);`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `ALTER TABLE users ADD COLUMN first_name TEXT NOT NULL DEFAULT '';`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `ALTER TABLE users ADD COLUMN last_name TEXT NOT NULL DEFAULT '';`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `ALTER TABLE users ADD COLUMN household_role TEXT NOT NULL DEFAULT 'member';`); err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, `SELECT id FROM users`)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var uid int64
		if err := rows.Scan(&uid); err != nil {
			return err
		}
		now := timeutil.NowSQLiteUTC()
		res, err := tx.ExecContext(ctx, `INSERT INTO households (name, created_at) VALUES ('My household', ?)`, now)
		if err != nil {
			return err
		}
		hid, err := res.LastInsertId()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE users SET household_id = ?, household_role = 'owner' WHERE id = ?`, hid, uid); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (3)`)
	return err
}

func migrateV4(ctx context.Context, tx *sql.Tx) error {
	var hasTZ bool
	err := tx.QueryRowContext(ctx, `
SELECT EXISTS (SELECT 1 FROM pragma_table_info('users') WHERE name = 'timezone')`).Scan(&hasTZ)
	if err != nil {
		return err
	}
	if hasTZ {
		if _, err := tx.ExecContext(ctx, `ALTER TABLE users DROP COLUMN timezone`); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (4)`)
	return err
}

func migrateV5(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `ALTER TABLE categories ADD COLUMN color TEXT NOT NULL DEFAULT ''`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (5)`)
	return err
}

func migrateV6(ctx context.Context, tx *sql.Tx) error {
	var hasUserID int
	if err := tx.QueryRowContext(ctx, `
SELECT COUNT(*) FROM pragma_table_info('categories') WHERE name = 'user_id'`).Scan(&hasUserID); err != nil {
		return err
	}
	if hasUserID == 0 {
		_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (6)`)
		return err
	}

	if _, err := tx.ExecContext(ctx, `ALTER TABLE categories ADD COLUMN household_id INTEGER REFERENCES households(id)`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
UPDATE categories SET household_id = (
  SELECT household_id FROM users WHERE users.id = categories.user_id
) WHERE household_id IS NULL`); err != nil {
		return err
	}
	if err := mergeCategoryDuplicatesByHousehold(ctx, tx); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
CREATE TABLE categories_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  household_id INTEGER NOT NULL REFERENCES households(id) ON DELETE CASCADE,
  name TEXT NOT NULL COLLATE NOCASE,
  icon TEXT NOT NULL DEFAULT '',
  color TEXT NOT NULL DEFAULT '',
  UNIQUE (household_id, name)
)`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO categories_new (id, household_id, name, icon, color)
SELECT id, household_id, name, icon, color FROM categories`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `PRAGMA foreign_keys=OFF`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DROP TABLE categories`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `ALTER TABLE categories_new RENAME TO categories`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `CREATE INDEX idx_categories_household ON categories(household_id)`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `PRAGMA foreign_keys=ON`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (6)`)
	return err
}

func mergeCategoryDuplicatesByHousehold(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(ctx, `SELECT id, household_id, name FROM categories ORDER BY id`)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	type row struct {
		id, hid int64
		name    string
	}
	var all []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.hid, &r.name); err != nil {
			return err
		}
		if r.hid == 0 {
			return fmt.Errorf("category %d has no household_id after backfill", r.id)
		}
		all = append(all, r)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	groups := make(map[string][]int64)
	for _, r := range all {
		key := fmt.Sprintf("%d|%s", r.hid, strings.ToLower(strings.TrimSpace(r.name)))
		groups[key] = append(groups[key], r.id)
	}
	for _, ids := range groups {
		if len(ids) < 2 {
			continue
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		survivor := ids[0]
		for _, dup := range ids[1:] {
			if _, err := tx.ExecContext(ctx, `UPDATE transactions SET category_id = ? WHERE category_id = ?`, survivor, dup); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `DELETE FROM categories WHERE id = ?`, dup); err != nil {
				return err
			}
		}
	}
	return nil
}
