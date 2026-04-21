package db

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func TestOpen_memoryDoesNotUseWALJournalMode(t *testing.T) {
	t.Parallel()
	d, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var jm string
	if err := d.QueryRowContext(context.Background(), `PRAGMA journal_mode`).Scan(&jm); err != nil {
		t.Fatal(err)
	}
	if strings.EqualFold(jm, "wal") {
		t.Fatalf("PRAGMA journal_mode = %q want non-wal (Open skips WAL pragma for :memory:)", jm)
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

func TestOpen_fileTuningPragmas(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "pragma.db")
	d, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var journal string
	if err := d.QueryRowContext(context.Background(), `PRAGMA journal_mode`).Scan(&journal); err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(journal, "wal") {
		t.Fatalf("PRAGMA journal_mode = %q want wal (file-backed DB must use WAL)", journal)
	}
	var cache int64
	if err := d.QueryRowContext(context.Background(), `PRAGMA cache_size`).Scan(&cache); err != nil {
		t.Fatal(err)
	}
	if cache != -16000 {
		t.Fatalf("PRAGMA cache_size = %d want -16000 (16 MiB page cache)", cache)
	}
	const wantMmap = int64(64 << 20)
	var mmap int64
	if err := d.QueryRowContext(context.Background(), `PRAGMA mmap_size`).Scan(&mmap); err != nil {
		t.Fatal(err)
	}
	if mmap != wantMmap {
		t.Fatalf("PRAGMA mmap_size = %d want %d", mmap, wantMmap)
	}
	var tempStore int
	if err := d.QueryRowContext(context.Background(), `PRAGMA temp_store`).Scan(&tempStore); err != nil {
		t.Fatal(err)
	}
	// SQLite: 2 = MEMORY (see PRAGMA temp_store).
	if tempStore != 2 {
		t.Fatalf("PRAGMA temp_store = %d want 2 (MEMORY)", tempStore)
	}
}

func TestOpen_notificationsTableAndIndex(t *testing.T) {
	t.Parallel()
	d, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var n int
	err = d.QueryRowContext(context.Background(), `
SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'notifications'`).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("notifications table: count %d want 1", n)
	}
	err = d.QueryRowContext(context.Background(), `
SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_notifications_user_created'`).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("idx_notifications_user_created: count %d want 1", n)
	}
}

func TestOpen_idxNotificationsUserUnreadPartialIndex(t *testing.T) {
	t.Parallel()
	d, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var n int
	err = d.QueryRowContext(context.Background(), `
SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_notifications_user_unread'`).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("idx_notifications_user_unread: count %d want 1", n)
	}
}

func TestOpen_idxTransactionsOccurredAt(t *testing.T) {
	t.Parallel()
	d, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var n int
	err = d.QueryRowContext(context.Background(), `
SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_transactions_occurred_at'`).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("idx_transactions_occurred_at: count %d want 1", n)
	}
}

func TestOpen_concurrentReadQueries(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "concurrent.db")
	d, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	const n = 32
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var one int
			errs <- d.QueryRowContext(context.Background(), `SELECT 1`).Scan(&one)
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}
