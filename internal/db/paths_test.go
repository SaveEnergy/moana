package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDBParentDir_createsNestedDirectories(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "nested", "deep", "moana.db")
	if err := ensureDBParentDir(filepath.Clean(path)); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(filepath.Join(root, "nested", "deep"))
	if err != nil {
		t.Fatal(err)
	}
	if !st.IsDir() {
		t.Fatal("expected directory")
	}
}

func TestEnsureDBParentDir_currentDirNoMkdir(t *testing.T) {
	t.Parallel()
	// filepath.Dir("foo.db") is "." — ensureDBParentDir returns early without error.
	if err := ensureDBParentDir("foo.db"); err != nil {
		t.Fatal(err)
	}
}
