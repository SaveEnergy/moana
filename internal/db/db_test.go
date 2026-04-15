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
