package handlers_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"moana/internal/testutil"
)

func TestLoginSubmit_wrongPasswordRendersForm(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "badpw@integration.test", "correct-secret", "user")
	client := testutil.NewCookieClient(t)
	resp, err := client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"badpw@integration.test"},
		"password": {"wrong-password"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 with login form", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Invalid email or password.") {
		t.Fatalf("expected auth error copy, got: %s", s[:min(500, len(s))])
	}
}

func TestLoginPage_sessionExpiredBanner(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	resp, err := http.Get(srv.URL + "/login?error=1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Session expired or invalid. Please sign in again.") {
		t.Fatalf("expected session banner, got: %s", s[:min(600, len(s))])
	}
}
