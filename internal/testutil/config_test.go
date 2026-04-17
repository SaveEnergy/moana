package testutil

import (
	"testing"
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
	if c.RequestTimeout <= 0 || c.SessionMaxAge <= 0 {
		t.Fatalf("timeouts must be positive: RequestTimeout=%v SessionMaxAge=%v", c.RequestTimeout, c.SessionMaxAge)
	}
}
