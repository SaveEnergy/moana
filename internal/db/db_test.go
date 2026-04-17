package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenMigrateFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	d, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var v int
	err = d.QueryRowContext(context.Background(), `SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != LatestMigrationVersion() {
		t.Fatalf("schema version = %d", v)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestOpenMemory(t *testing.T) {
	t.Parallel()
	d, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var v int
	err = d.QueryRowContext(context.Background(), `SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != LatestMigrationVersion() {
		t.Fatalf("schema version = %d", v)
	}
}

func TestOpen_secondPingSucceeds(t *testing.T) {
	t.Parallel()
	d, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	// Open already pings after migrate; pool must remain usable.
	if err := d.PingContext(context.Background()); err != nil {
		t.Fatalf("post-Open ping: %v", err)
	}
}

func TestOpen_synchronousNormal(t *testing.T) {
	t.Parallel()
	d, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var sync int
	err = d.QueryRowContext(context.Background(), `PRAGMA synchronous`).Scan(&sync)
	if err != nil {
		t.Fatal(err)
	}
	// SQLite: 1 = NORMAL (set in Open after WAL/busy_timeout).
	if sync != 1 {
		t.Fatalf("PRAGMA synchronous = %d want 1 (NORMAL)", sync)
	}
}
