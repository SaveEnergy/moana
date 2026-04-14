package dbutil

import (
	"database/sql"

	"moana/internal/db"
	"moana/internal/store"
)

// OpenStore opens SQLite at path and returns a [store.Store] and the underlying [*sql.DB] (caller must Close the DB).
func OpenStore(path string) (*store.Store, *sql.DB, error) {
	database, err := db.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return store.New(database), database, nil
}
