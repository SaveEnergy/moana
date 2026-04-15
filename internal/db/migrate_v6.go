package db

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

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
