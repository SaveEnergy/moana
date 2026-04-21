package db

import (
	"context"
	"database/sql"
)

// migrateV10 adds a partial index for unread notifications (read_at IS NULL) per user.
// Speeds up COUNT/list filters for unread-only workloads (e.g. future badge queries).
func migrateV10(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `
CREATE INDEX idx_notifications_user_unread ON notifications(user_id) WHERE read_at IS NULL`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (10)`)
	return err
}
