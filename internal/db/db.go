package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	_ "modernc.org/sqlite"
)

var memoryDBSeq atomic.Uint64

// Open opens the SQLite database with WAL and foreign keys enabled.
// Use path ":memory:" for an in-memory database (tests and ephemeral runs).
func Open(path string) (*sql.DB, error) {
	var dsn string
	switch path {
	case ":memory:":
		// Unique URI per Open so parallel tests do not share one DB.
		id := memoryDBSeq.Add(1)
		dsn = fmt.Sprintf("file:memdb%d?mode=memory&cache=shared&_pragma=foreign_keys(1)", id)
	default:
		clean := filepath.Clean(path)
		if err := ensureDBParentDir(clean); err != nil {
			return nil, err
		}
		dsn = "file:" + strings.ReplaceAll(clean, "\\", "/") + "?cache=shared&_pragma=foreign_keys(1)"
	}
	d, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(1)

	if path != ":memory:" {
		if _, err := d.ExecContext(context.Background(), `PRAGMA journal_mode = WAL;`); err != nil {
			_ = d.Close()
			return nil, fmt.Errorf("wal: %w", err)
		}
	}
	if _, err := d.ExecContext(context.Background(), `PRAGMA busy_timeout = 5000;`); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("busy_timeout: %w", err)
	}

	if err := migrate(d); err != nil {
		_ = d.Close()
		return nil, err
	}
	return d, nil
}

func migrate(d *sql.DB) error {
	ctx := context.Background()
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_version (
  version INTEGER NOT NULL PRIMARY KEY
);`); err != nil {
		return err
	}

	var v int
	err = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&v)
	if err != nil {
		return err
	}

	if v < 1 {
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
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (1)`); err != nil {
			return err
		}
		v = 1
	}

	if v < 2 {
		if _, err := tx.ExecContext(ctx, `ALTER TABLE categories ADD COLUMN icon TEXT NOT NULL DEFAULT '';`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (2)`); err != nil {
			return err
		}
		v = 2
	}

	if v < 3 {
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
			now := time.Now().UTC().Format(time.RFC3339Nano)
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
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (3)`); err != nil {
			return err
		}
		v = 3
	}

	if v < 4 {
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
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (4)`); err != nil {
			return err
		}
		v = 4
	}

	if v < 5 {
		if _, err := tx.ExecContext(ctx, `ALTER TABLE categories ADD COLUMN color TEXT NOT NULL DEFAULT ''`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (5)`); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func ensureDBParentDir(cleanPath string) error {
	dir := filepath.Dir(cleanPath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
