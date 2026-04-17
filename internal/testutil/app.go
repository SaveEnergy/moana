package testutil

import (
	"database/sql"
	"testing"
	"time"

	moanaapp "moana/internal/app"
	"moana/internal/config"
	"moana/internal/dbutil"
	"moana/internal/handlers"
)

// DefaultTestRepoURL matches [config.Load] when MOANA_REPO_URL is unset (footer + login GitHub link).
const DefaultTestRepoURL = "https://github.com/SaveEnergy/moana"

// DefaultTestConfig returns a config suitable for integration tests (in-memory DB, non-secure cookies).
func DefaultTestConfig() *config.Config {
	return &config.Config{
		Listen:         ":0",
		DBPath:         ":memory:",
		SessionSecret:  []byte("integration-test-session-secret-32b!"),
		SecureCookies:  false,
		SessionMaxAge:  time.Hour,
		RequestTimeout: 30 * time.Second,
		RepoURL:        DefaultTestRepoURL,
	}
}

// NewAppWithDB returns a handlers.App, the underlying [*sql.DB] (for tests that need to close
// or inspect the connection), and a cleanup that closes the database.
func NewAppWithDB(t *testing.T) (*handlers.App, *sql.DB, func()) {
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
	return a, sqlDB, func() { sqlDB.Close() }
}

// NewApp returns a handlers.App backed by SQLite :memory: and a cleanup that closes the DB.
func NewApp(t *testing.T) (*handlers.App, func()) {
	app, _, cleanup := NewAppWithDB(t)
	return app, cleanup
}
