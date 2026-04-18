package handlers

import (
	"testing"

	"moana/internal/safepath"
)

func TestLoginRedirectAuth_matchesLoginPathWithErrorQuery(t *testing.T) {
	t.Parallel()
	if want := loginPath + "?error=1"; loginRedirectAuth != want {
		t.Fatalf("loginRedirectAuth=%q want %q (middleware WithAuth must match [routes_auth] GET %s)",
			loginRedirectAuth, want, loginPath)
	}
}

func TestHistoryPath_matchesSafepathDefault(t *testing.T) {
	t.Parallel()
	if historyPath != safepath.Default {
		t.Fatalf("historyPath=%q must equal safepath.Default (%q) so ledger redirects match safe internal paths",
			historyPath, safepath.Default)
	}
}
