package handlers

import (
	"testing"

	"moana/internal/safepath"
)

func TestLoginRedirectAuth_matchesLoginPathWithErrorQuery(t *testing.T) {
	t.Parallel()
	if want := LoginPath + "?error=1"; LoginRedirectAuth != want {
		t.Fatalf("LoginRedirectAuth=%q want %q (middleware WithAuth must match [routes_auth] GET %s)",
			LoginRedirectAuth, want, LoginPath)
	}
}

func TestHistoryPath_matchesSafepathDefault(t *testing.T) {
	t.Parallel()
	if HistoryPath != safepath.Default {
		t.Fatalf("HistoryPath=%q must equal safepath.Default (%q) so ledger redirects match safe internal paths",
			HistoryPath, safepath.Default)
	}
}
