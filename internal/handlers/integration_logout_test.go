package handlers_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"moana/internal/testutil"
)

func TestLogout_redirectsToLoginAndClearsSession(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "logout@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "logout@integration.test", "pw")

	resp, err := client.PostForm(srv.URL+"/logout", url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("after logout redirect chain status %d want 200", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Sign in to your account") {
		t.Fatal("expected login page after logout")
	}

	resp2, err := client.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	if !strings.Contains(string(body2), "Sign in to your account") {
		t.Fatal("expected auth gate to login after session cleared")
	}
}

func TestLogout_withoutSessionRedirectsToLogin(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm(srv.URL+"/logout", url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("status %d want %d", resp.StatusCode, http.StatusSeeOther)
	}
	loc := resp.Header.Get("Location")
	if !strings.Contains(loc, "/login") {
		t.Fatalf("unexpected Location %q", loc)
	}
}
