package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"sync/atomic"

	_ "modernc.org/sqlite"
)

var memoryDBSeq atomic.Uint64

// Open opens the SQLite database with WAL and foreign keys enabled.
// Use path ":memory:" for an in-memory database (tests and ephemeral runs).
func Open(path string) (*sql.DB, error) {
	var dsn string
	switch path {
	case ":memory:":
		// Unique URI per Open so parallel tests do not share one DB.
		id := memoryDBSeq.Add(1)
		dsn = fmt.Sprintf("file:memdb%d?mode=memory&cache=shared&_pragma=foreign_keys(1)", id)
	default:
		clean := filepath.Clean(path)
		if err := ensureDBParentDir(clean); err != nil {
			return nil, err
		}
		dsn = "file:" + strings.ReplaceAll(clean, "\\", "/") + "?cache=shared&_pragma=foreign_keys(1)"
	}
	d, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(1)

	if path != ":memory:" {
		if _, err := d.ExecContext(context.Background(), `PRAGMA journal_mode = WAL;`); err != nil {
			_ = d.Close()
			return nil, fmt.Errorf("wal: %w", err)
		}
	}
	if _, err := d.ExecContext(context.Background(), `PRAGMA busy_timeout = 5000;`); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("busy_timeout: %w", err)
	}

	if err := migrate(d); err != nil {
		_ = d.Close()
		return nil, err
	}
	return d, nil
}
