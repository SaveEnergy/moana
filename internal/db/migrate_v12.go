package db

import (
	"context"
	"database/sql"
)

// migrateV12 adds per-user profile avatar version (0 = use initials; >0 = JPEG on disk, cache-bust via query).
func migrateV12(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `ALTER TABLE users ADD COLUMN avatar_rev INTEGER NOT NULL DEFAULT 0`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (12)`)
	return err
}
