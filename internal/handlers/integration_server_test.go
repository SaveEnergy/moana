package handlers_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"moana/internal/testutil"
)

func TestHealth(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	srv := testutil.NewServer(t, app)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "ok" {
		t.Fatalf("body %q", b)
	}
}

func TestLoginAndOverview(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "user@integration.test", "correct-password", "user")

	srv := testutil.NewServer(t, app)
	defer srv.Close()

	client := testutil.NewCookieClient(t)

	resp, err := client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"user@integration.test"},
		"password": {"correct-password"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `class="dashboard-hero-title"`) {
		snippet := string(body)
		if len(snippet) > 300 {
			snippet = snippet[:300] + "…"
		}
		t.Fatalf("expected dashboard hero, got: %s", snippet)
	}
}

func TestSettingsPageOKForLoggedInUser(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "plain@moana.test", "pw", "user")
	srv := testutil.NewServer(t, app)
	defer srv.Close()
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "plain@moana.test", "pw")
	resp, err := client.Get(srv.URL + "/settings")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "settings-page") {
		t.Fatalf("expected settings page")
	}
}
