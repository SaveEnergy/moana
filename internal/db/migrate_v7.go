package db

import (
	"context"
	"database/sql"
)

// migrateV7 adds indexes for common household-scoped lookups (categories already have idx_categories_household from v6).
func migrateV7(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_users_household_id ON users(household_id)`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (7)`)
	return err
}
