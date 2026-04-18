package handlers

import (
	"net/http"
	"strings"
	"testing"

	"moana/internal/safepath"
)

func TestDashboardRootPattern_exactRootServeMuxString(t *testing.T) {
	t.Parallel()
	const want = "GET /{$}"
	if got := http.MethodGet + " " + DashboardRootPattern; got != want {
		t.Fatalf("dashboard mux pattern %q must stay %q (stdlib root exact-match)", got, want)
	}
}

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

func TestTransactionPathPatterns_useStdlibWildcardSegment(t *testing.T) {
	t.Parallel()
	for _, p := range []string{TransactionPathPattern, TransactionEditPathPattern} {
		if !strings.Contains(p, "{id}") {
			t.Fatalf("%q must contain {id} (Go 1.22+ ServeMux named wildcard)", p)
		}
	}
}

func TestCategoryMutatePaths_haveCategoriesPrefix(t *testing.T) {
	t.Parallel()
	for _, p := range []string{CategoriesUpdatePath, CategoriesDeletePath} {
		if !strings.HasPrefix(p, CategoriesPath+"/") {
			t.Fatalf("%q must be under %q/", p, CategoriesPath)
		}
	}
}

func TestSettingsPOSTPaths_haveSettingsPrefix(t *testing.T) {
	t.Parallel()
	for _, p := range []string{
		SettingsProfilePath,
		SettingsHouseholdPath,
		SettingsHouseholdMembersPath,
		SettingsHouseholdMembersRemovePath,
	} {
		if !strings.HasPrefix(p, SettingsPath+"/") {
			t.Fatalf("%q must be under %q/", p, SettingsPath)
		}
	}
}
