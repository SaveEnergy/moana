package db

import (
	"context"
	"database/sql"
)

// migrateV9 adds per-user notifications (in-app inbox).
func migrateV9(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `
CREATE TABLE notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  body TEXT NOT NULL,
  read_at TEXT,
  created_at TEXT NOT NULL
);`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at)`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (9)`)
	return err
}
