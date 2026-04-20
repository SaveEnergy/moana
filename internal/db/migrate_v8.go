package db

import (
	"context"
	"database/sql"
)

// migrateV8 adds an index on transactions.occurred_at for time-range filters and
// ORDER BY occurred_at (listing, history, aggregates) combined with joins to users.
func migrateV8(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_transactions_occurred_at ON transactions(occurred_at)`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (8)`)
	return err
}
