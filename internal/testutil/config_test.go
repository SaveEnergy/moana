package testutil

import (
	"testing"
	"time"
)

func TestDefaultTestConfig_repoURLMatchesDefaultRepoConstant(t *testing.T) {
	t.Parallel()
	c := DefaultTestConfig()
	if c.RepoURL != DefaultTestRepoURL {
		t.Fatalf("RepoURL %q want %q", c.RepoURL, DefaultTestRepoURL)
	}
}

func TestDefaultTestConfig_integrationDefaults(t *testing.T) {
	t.Parallel()
	c := DefaultTestConfig()
	if c.DBPath != ":memory:" {
		t.Fatalf("DBPath %q want :memory:", c.DBPath)
	}
	if len(c.SessionSecret) < 32 {
		t.Fatalf("SessionSecret length %d (want enough bytes for HMAC)", len(c.SessionSecret))
	}
	if c.SecureCookies {
		t.Fatal("SecureCookies must be false for httptest clients without TLS")
	}
	if c.RequestTimeout != 30*time.Second {
		t.Fatalf("RequestTimeout %v want 30s (faster than [config.Load] default for integration tests)", c.RequestTimeout)
	}
	if c.SessionMaxAge != time.Hour {
		t.Fatalf("SessionMaxAge %v want 1h", c.SessionMaxAge)
	}
}
