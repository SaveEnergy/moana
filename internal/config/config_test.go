package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoad_productionRequiresSessionSecret(t *testing.T) {
	t.Setenv("MOANA_ENV", "production")
	t.Setenv("MOANA_SESSION_SECRET", "")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when MOANA_ENV=production and secret empty")
	}
	if !strings.Contains(err.Error(), "MOANA_SESSION_SECRET") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestLoad_developmentAllowsEmptySecret(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.SessionSecret) == 0 {
		t.Fatal("expected dev fallback secret bytes")
	}
	if c.SecureCookies {
		t.Fatal("SecureCookies must be false in development")
	}
}

func TestLoad_productionSetsSecureCookies(t *testing.T) {
	t.Setenv("MOANA_ENV", "production")
	t.Setenv("MOANA_SESSION_SECRET", "minimum-production-secret-for-test")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !c.SecureCookies {
		t.Fatal("SecureCookies must be true in production")
	}
}

func TestLoad_productionRejectsShortSessionSecret(t *testing.T) {
	t.Setenv("MOANA_ENV", "production")
	t.Setenv("MOANA_SESSION_SECRET", "too-short-for-production-use")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when MOANA_SESSION_SECRET is shorter than 32 chars in production")
	}
	if !strings.Contains(err.Error(), "at least") || !strings.Contains(err.Error(), "32") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestLoad_defaultRepoURL(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_REPO_URL", "")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	const want = "https://github.com/SaveEnergy/moana"
	if c.RepoURL != want {
		t.Fatalf("RepoURL %q want %q", c.RepoURL, want)
	}
}

func TestLoad_repoURLFromEnv(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_REPO_URL", "https://example.com/my-fork")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.RepoURL != "https://example.com/my-fork" {
		t.Fatalf("RepoURL %q", c.RepoURL)
	}
}

func TestLoad_defaultRequestTimeout(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_REQUEST_TIMEOUT_SEC", "")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.RequestTimeout != 60*time.Second {
		t.Fatalf("RequestTimeout %v want 60s", c.RequestTimeout)
	}
}

func TestLoad_requestTimeoutFromEnv(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_REQUEST_TIMEOUT_SEC", "42")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.RequestTimeout != 42*time.Second {
		t.Fatalf("RequestTimeout %v want 42s", c.RequestTimeout)
	}
}

func TestLoad_invalidRequestTimeoutFallsBackToDefault(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_REQUEST_TIMEOUT_SEC", "not-a-number")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.RequestTimeout != 60*time.Second {
		t.Fatalf("RequestTimeout %v want 60s fallback", c.RequestTimeout)
	}
}

func TestLoad_sessionMaxAgeFromEnv(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_SESSION_MAX_AGE_SEC", "3600")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.SessionMaxAge != time.Hour {
		t.Fatalf("SessionMaxAge %v want 1h", c.SessionMaxAge)
	}
}

func TestLoad_invalidSessionMaxAgeFallsBackToDefault(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_SESSION_MAX_AGE_SEC", "not-a-number")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := 604800 * time.Second // 7d, matches Load default when Atoi fails or <= 0
	if c.SessionMaxAge != want {
		t.Fatalf("SessionMaxAge %v want %v", c.SessionMaxAge, want)
	}
}

func TestDBPath_defaultWhenUnsetOrEmpty(t *testing.T) {
	t.Setenv("MOANA_DB_PATH", "")
	if got := DBPath(); got != "data/moana.db" {
		t.Fatalf("DBPath() = %q want data/moana.db", got)
	}
}

func TestDBPath_fromEnv(t *testing.T) {
	t.Setenv("MOANA_DB_PATH", "/tmp/moana-dbpath-regression.db")
	if got := DBPath(); got != "/tmp/moana-dbpath-regression.db" {
		t.Fatalf("DBPath() = %q", got)
	}
}

func TestLoad_listenAndDbPathFromEnv(t *testing.T) {
	t.Setenv("MOANA_ENV", "development")
	t.Setenv("MOANA_SESSION_SECRET", "")
	t.Setenv("MOANA_LISTEN", ":9999")
	t.Setenv("MOANA_DB_PATH", "/tmp/moana-config-test.db")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.Listen != ":9999" {
		t.Fatalf("Listen %q want :9999", c.Listen)
	}
	if c.DBPath != "/tmp/moana-config-test.db" {
		t.Fatalf("DBPath %q", c.DBPath)
	}
}
