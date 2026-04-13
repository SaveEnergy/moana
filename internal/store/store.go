package store

import "database/sql"

// Store wraps database access for the app.
type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}
