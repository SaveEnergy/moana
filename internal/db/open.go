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

// MaxOpenConns is the [database/sql] pool size for SQLite opened by [Open]. Under WAL, readers can
// overlap while writes stay serialized; a pool of one would serialize concurrent read workloads
// (e.g. [moana/internal/dashboard.BuildPageData] overlapping aggregate + heatmap + recent queries).
// Keep this greater than 1 so those reads can run concurrently.
const MaxOpenConns = 8

var memoryDBSeq atomic.Uint64

// sqliteDSN builds a modernc.org/sqlite driver URI. Each ":memory:" open uses a unique name so
// parallel tests do not share one database. File paths use forward slashes in the URI.
func sqliteDSN(path string) (string, error) {
	switch path {
	case ":memory:":
		id := memoryDBSeq.Add(1)
		return fmt.Sprintf("file:memdb%d?mode=memory&cache=shared&_pragma=foreign_keys(1)", id), nil
	default:
		clean := filepath.Clean(path)
		if err := ensureDBParentDir(clean); err != nil {
			return "", err
		}
		return "file:" + strings.ReplaceAll(clean, "\\", "/") + "?cache=shared&_pragma=foreign_keys(1)", nil
	}
}

// Open opens the SQLite database with WAL and foreign keys enabled.
// Use path ":memory:" for an in-memory database (tests and ephemeral runs).
func Open(path string) (*sql.DB, error) {
	dsn, err := sqliteDSN(path)
	if err != nil {
		return nil, err
	}
	d, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(MaxOpenConns)
	d.SetMaxIdleConns(MaxOpenConns)

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
	// NORMAL balances durability with throughput under WAL (default is often FULL).
	if _, err := d.ExecContext(context.Background(), `PRAGMA synchronous = NORMAL;`); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("synchronous: %w", err)
	}

	if path != ":memory:" {
		// Larger cache + mmap reduce read I/O vs SQLite defaults (small cache, no mmap).
		if _, err := d.ExecContext(context.Background(), `PRAGMA cache_size = -16000;`); err != nil {
			_ = d.Close()
			return nil, fmt.Errorf("cache_size: %w", err)
		}
		const mmapBytes = 64 << 20 // 64 MiB upper bound for memory-mapped reads
		if _, err := d.ExecContext(context.Background(), fmt.Sprintf("PRAGMA mmap_size = %d", mmapBytes)); err != nil {
			_ = d.Close()
			return nil, fmt.Errorf("mmap_size: %w", err)
		}
		if _, err := d.ExecContext(context.Background(), `PRAGMA temp_store = MEMORY;`); err != nil {
			_ = d.Close()
			return nil, fmt.Errorf("temp_store: %w", err)
		}
	}

	if err := migrate(d); err != nil {
		_ = d.Close()
		return nil, err
	}
	if err := d.PingContext(context.Background()); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return d, nil
}
