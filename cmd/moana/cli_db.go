package main

import (
	"database/sql"

	"moana/internal/config"
	"moana/internal/dbutil"
	"moana/internal/store"
)

// openCLIStore opens the database using MOANA_DB_PATH (same default as [config.DBPath]).
func openCLIStore() (*store.Store, *sql.DB, error) {
	return dbutil.OpenStore(config.DBPath())
}
