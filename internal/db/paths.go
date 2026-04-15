package db

import (
	"os"
	"path/filepath"
)

func ensureDBParentDir(cleanPath string) error {
	dir := filepath.Dir(cleanPath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
