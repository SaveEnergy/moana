package handlers

import (
	"strings"
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

func TestCategoryPaths_derivedFromCategoriesPath(t *testing.T) {
	t.Parallel()
	if got, want := CategoriesUpdatePath, CategoriesPath+"/update"; got != want {
		t.Fatalf("CategoriesUpdatePath=%q want %q", got, want)
	}
	if got, want := CategoriesDeletePath, CategoriesPath+"/delete"; got != want {
		t.Fatalf("CategoriesDeletePath=%q want %q", got, want)
	}
}

func TestSettingsPaths_derivedFromSettingsPath(t *testing.T) {
	t.Parallel()
	if got, want := SettingsProfilePath, SettingsPath+"/profile"; got != want {
		t.Fatalf("SettingsProfilePath=%q want %q", got, want)
	}
	if got, want := SettingsHouseholdPath, SettingsPath+"/household"; got != want {
		t.Fatalf("SettingsHouseholdPath=%q want %q", got, want)
	}
	if got, want := SettingsHouseholdMembersPath, SettingsHouseholdPath+"/members"; got != want {
		t.Fatalf("SettingsHouseholdMembersPath=%q want %q", got, want)
	}
	if got, want := SettingsHouseholdMembersRemovePath, SettingsHouseholdMembersPath+"/remove"; got != want {
		t.Fatalf("SettingsHouseholdMembersRemovePath=%q want %q", got, want)
	}
}
