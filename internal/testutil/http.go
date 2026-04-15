package testutil

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"moana/internal/handlers"
	"moana/internal/server"
)

// NewServer returns an httptest.Server with routes and logging disabled.
func NewServer(t *testing.T, app *handlers.App) *httptest.Server {
	t.Helper()
	return httptest.NewServer(server.NewRouterWithRouterOptions(&server.RouterOptions{DisableRequestLogging: true}, app))
}

// NewAppServer combines [NewApp] and [NewServer]; cleanup closes the server then the database.
func NewAppServer(t *testing.T) (*handlers.App, *httptest.Server, func()) {
	t.Helper()
	app, cleanup := NewApp(t)
	srv := NewServer(t, app)
	return app, srv, func() {
		srv.Close()
		cleanup()
	}
}

// NewCookieClient returns an [http.Client] with a cookie jar (for session cookies in tests).
func NewCookieClient(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Client{Jar: jar}
}

// MustLogin POSTs to baseURL/login with the given credentials. The HTTP client follows redirects;
// the final response after a successful login is typically 200 from GET /.
// baseURL must be the server root (no trailing slash), e.g. httptest.Server.URL.
func MustLogin(t *testing.T, client *http.Client, baseURL, email, password string) {
	t.Helper()
	resp, err := client.PostForm(baseURL+"/login", url.Values{
		"email":    {email},
		"password": {password},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: status %d", resp.StatusCode)
	}
}
