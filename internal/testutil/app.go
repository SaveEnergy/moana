package testutil

import (
	"testing"
	"time"

	moanaapp "moana/internal/app"
	"moana/internal/config"
	"moana/internal/dbutil"
	"moana/internal/handlers"
)

// DefaultTestConfig returns a config suitable for integration tests (in-memory DB, non-secure cookies).
func DefaultTestConfig() *config.Config {
	return &config.Config{
		Listen:         ":0",
		DBPath:         ":memory:",
		SessionSecret:  []byte("integration-test-session-secret-32b!"),
		SecureCookies:  false,
		SessionMaxAge:  time.Hour,
		RequestTimeout: 30 * time.Second,
	}
}

// NewApp returns a handlers.App backed by SQLite :memory: and a cleanup that closes the DB.
func NewApp(t *testing.T) (*handlers.App, func()) {
	t.Helper()
	st, sqlDB, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	cfg := DefaultTestConfig()
	a, err := moanaapp.New(cfg, st)
	if err != nil {
		sqlDB.Close()
		t.Fatal(err)
	}
	return a, func() { sqlDB.Close() }
}
