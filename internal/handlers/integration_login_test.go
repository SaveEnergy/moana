package handlers_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"moana/internal/testutil"
)

func TestLoginSubmit_trimmedEmailAuthenticates(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "trimlogin@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	resp, err := client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"  trimlogin@integration.test  "},
		"password": {"pw"},
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
