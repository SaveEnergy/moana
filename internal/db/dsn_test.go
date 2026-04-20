package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSqliteDSN_memoryUniqueAndSharedCache(t *testing.T) {
	t.Parallel()
	a, err := sqliteDSN(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	b, err := sqliteDSN(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("expected distinct in-memory DSNs for parallel isolation")
	}
	for _, s := range []string{a, b} {
		if !strings.Contains(s, "mode=memory") || !strings.Contains(s, "cache=shared") || !strings.Contains(s, "foreign_keys(1)") {
			t.Fatalf("dsn %q missing expected fragments", s)
		}
	}
}

func TestSqliteDSN_fileUsesForwardSlashes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "nested", "app.db")
	dsn, err := sqliteDSN(p)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(dsn, `\`) {
		t.Fatalf("dsn %q must not contain backslashes (SQLite URI)", dsn)
	}
	if !strings.HasPrefix(dsn, "file:") || !strings.Contains(dsn, "nested/app.db") {
		t.Fatalf("dsn %q", dsn)
	}
}
