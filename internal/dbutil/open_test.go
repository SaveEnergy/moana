package dbutil

import (
	"context"
	"path/filepath"
	"testing"

	"moana/internal/passwordtest"
)

func TestOpenStore_memory(t *testing.T) {
	t.Parallel()
	st, db, err := OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if st == nil {
		t.Fatal("nil store")
	}
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	if _, err := st.CreateUser(ctx, "openstore@example.com", hash, "user"); err != nil {
		t.Fatal(err)
	}
}

func TestOpenStore_filePath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "openstore-file.db")
	st, db, err := OpenStore(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if st == nil {
		t.Fatal("nil store")
	}
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	if _, err := st.CreateUser(ctx, "openstore-file@example.com", hash, "user"); err != nil {
		t.Fatal(err)
	}
}
