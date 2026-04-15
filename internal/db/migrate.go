package db

import (
	"context"
	"database/sql"
)

// migrationSteps runs in order; each migrateVN inserts schema_version (see migrate_steps.go, migrate_v6.go).
var migrationSteps = []func(context.Context, *sql.Tx) error{
	migrateV1,
	migrateV2,
	migrateV3,
	migrateV4,
	migrateV5,
	migrateV6,
	migrateV7,
}

// LatestMigrationVersion is the schema_version after [Open] runs all steps (len(migrationSteps)).
func LatestMigrationVersion() int {
	return len(migrationSteps)
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

	for i := range migrationSteps {
		want := i + 1
		if v < want {
			if err := migrationSteps[i](ctx, tx); err != nil {
				return err
			}
			v = want
		}
	}

	return tx.Commit()
}
