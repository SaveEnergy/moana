package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moana/internal/config"
)

func resolveAvatarDir(cfg *config.Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("nil config")
	}
	if d := strings.TrimSpace(cfg.AvatarDataDir); d != "" {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return "", err
		}
		return d, nil
	}
	if cfg.DBPath == ":memory:" {
		dir, err := os.MkdirTemp("", "moana-avatar-*")
		if err != nil {
			return "", err
		}
		if err := os.Chmod(dir, 0o700); err != nil {
			_ = os.RemoveAll(dir)
			return "", err
		}
		return dir, nil
	}
	dir := filepath.Join(filepath.Dir(cfg.DBPath), "avatars")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}
