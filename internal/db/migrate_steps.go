package db

import (
	"context"
	"database/sql"

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
