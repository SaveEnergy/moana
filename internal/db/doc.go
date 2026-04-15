// Package db opens SQLite (open.go + paths.go), applies migrations (migrate.go runs migrationSteps: v1–v5 in migrate_steps.go, v6 in migrate_v6.go, v7 in migrate_v7.go; [LatestMigrationVersion] matches applied steps), and exposes the shared [*sql.DB] used by [moana/internal/store].
package db
