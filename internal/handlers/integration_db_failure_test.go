package handlers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"moana/internal/handlers"
	"moana/internal/httperr"
	"moana/internal/testutil"
)

func TestCurrentUser_noSessionReturnsErrAuthRequired(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.NewApp(t)
	defer cleanup()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := app.CurrentUser(req)
	if !errors.Is(err, handlers.ErrAuthRequired) {
		t.Fatalf("got %v want %v", err, handlers.ErrAuthRequired)
	}
}

func TestAuthenticatedRoute_returns500WhenDatabaseClosed(t *testing.T) {
	t.Parallel()
	app, db, cleanup := testutil.NewAppWithDB(t)
	defer cleanup()
	srv := testutil.NewServer(t, app)
	defer srv.Close()

	testutil.MustCreateUser(t, app, "dbfail@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "dbfail@integration.test", "pw")

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	resp, err := client.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status %d want 500", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), httperr.InternalMessage) {
		t.Fatalf("body %q", body)
	}
}

func TestLoginPage_okWhenDatabaseClosedWithoutSession(t *testing.T) {
	t.Parallel()
	app, db, cleanup := testutil.NewAppWithDB(t)
	defer cleanup()
	srv := testutil.NewServer(t, app)
	defer srv.Close()
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	resp, err := http.Get(srv.URL + "/login")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Sign in to your account") {
		t.Fatalf("expected login page HTML")
	}
}
