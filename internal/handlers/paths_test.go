package handlers

import (
	"net/http"
	"net/url"
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
	if want := LoginPath + "?" + LoginErrorQueryParam + "=1"; LoginRedirectAuth != want {
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

func TestSettingsFlashQueryParams_roundtripInURLs(t *testing.T) {
	t.Parallel()
	errU, err := url.Parse(SettingsPath + "?" + SettingsErrorQueryParam + "=" + url.QueryEscape(`bad & "msg"`))
	if err != nil {
		t.Fatal(err)
	}
	if got := errU.Query().Get(SettingsErrorQueryParam); got != `bad & "msg"` {
		t.Fatalf("err flash: got %q", got)
	}
	okU, err := url.Parse(SettingsPath + "?" + SettingsOKQueryParam + "=" + url.QueryEscape("saved&key"))
	if err != nil {
		t.Fatal(err)
	}
	if got := okU.Query().Get(SettingsOKQueryParam); got != "saved&key" {
		t.Fatalf("ok flash: got %q", got)
	}
}

func TestTransactionPathPatterns_useStdlibWildcardSegment(t *testing.T) {
	t.Parallel()
	// Keep in sync with routes_ledger.go registrations (derived from TransactionsPath + segment).
	if got, want := TransactionPathPattern, TransactionsPath+"/{id}"; got != want {
		t.Fatalf("TransactionPathPattern=%q want %q", got, want)
	}
	if got, want := TransactionEditPathPattern, TransactionsPath+"/{id}/edit"; got != want {
		t.Fatalf("TransactionEditPathPattern=%q want %q", got, want)
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
