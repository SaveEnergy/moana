package handlers_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"moana/internal/handlers"
	"moana/internal/testutil"
)

func TestLoginSubmit_trimmedEmailAuthenticates(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "trimlogin@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	resp, err := client.PostForm(srv.URL+handlers.LoginPath, url.Values{
		handlers.LoginFieldEmail:    {"  trimlogin@integration.test  "},
		handlers.LoginFieldPassword: {"pw"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: status %d want 200", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `class="dashboard-hero-title"`) {
		t.Fatalf("expected dashboard after login with padded email; got body prefix %q", string(body[:min(400, len(body))]))
	}
}

func TestLoginSubmit_wrongPasswordShowsGenericError(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "badpw@integration.test", "correct-secret", "user")
	client := testutil.NewCookieClient(t)
	resp, err := client.PostForm(srv.URL+handlers.LoginPath, url.Values{
		handlers.LoginFieldEmail:    {"badpw@integration.test"},
		handlers.LoginFieldPassword: {"wrong-password"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: status %d want 200", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	const want = "Invalid email or password."
	if !strings.Contains(string(body), want) {
		t.Fatalf("expected error message in body; want substring %q; body prefix %q", want, string(body[:min(500, len(body))]))
	}
}

func TestLoginSubmit_unknownEmailShowsSameGenericError(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	client := testutil.NewCookieClient(t)
	resp, err := client.PostForm(srv.URL+handlers.LoginPath, url.Values{
		handlers.LoginFieldEmail:    {"nobody@integration.test"},
		handlers.LoginFieldPassword: {"any-password"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: status %d want 200", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	const want = "Invalid email or password."
	if !strings.Contains(string(body), want) {
		t.Fatalf("expected error message in body; want substring %q; body prefix %q", want, string(body[:min(500, len(body))]))
	}
}

func TestLoginPage_redirectsWhenAlreadyAuthenticated(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "logpage-redir@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "logpage-redir@integration.test", "pw")

	resp, err := client.Get(srv.URL + handlers.LoginPath)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 after redirect chain", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `class="dashboard-hero-title"`) {
		t.Fatalf("expected dashboard after GET /login with session; body prefix %q", string(body[:min(400, len(body))]))
	}
}
