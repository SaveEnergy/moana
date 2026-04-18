package handlers

import "testing"

func TestLoginRedirectAuth_matchesLoginPathWithErrorQuery(t *testing.T) {
	t.Parallel()
	if want := loginPath + "?error=1"; loginRedirectAuth != want {
		t.Fatalf("loginRedirectAuth=%q want %q (middleware WithAuth must match [routes_auth] GET %s)",
			loginRedirectAuth, want, loginPath)
	}
}
