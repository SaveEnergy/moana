package db

import (
	"context"
	"database/sql"
)

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
		if err := migrateV1(ctx, tx); err != nil {
			return err
		}
		v = 1
	}

	if v < 2 {
		if err := migrateV2(ctx, tx); err != nil {
			return err
		}
		v = 2
	}

	if v < 3 {
		if err := migrateV3(ctx, tx); err != nil {
			return err
		}
		v = 3
	}

	if v < 4 {
		if err := migrateV4(ctx, tx); err != nil {
			return err
		}
		v = 4
	}

	if v < 5 {
		if err := migrateV5(ctx, tx); err != nil {
			return err
		}
	}

	if v < 6 {
		if err := migrateV6(ctx, tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}
