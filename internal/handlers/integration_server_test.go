package handlers_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"moana/internal/testutil"
)

func TestLoginPageOK(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	resp, err := http.Get(srv.URL + "/login")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Sign in to your account") {
		t.Fatalf("expected login template")
	}
	if !strings.Contains(s, `href="`+testutil.DefaultTestRepoURL+`"`) {
		t.Fatalf("expected login footer repo link (set RepoURL in testutil.DefaultTestConfig)")
	}
}

func TestUnauthenticatedDashboardRedirectsToLogin(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("status %d want %d", resp.StatusCode, http.StatusSeeOther)
	}
	loc := resp.Header.Get("Location")
	if !strings.Contains(loc, "/login") || !strings.Contains(loc, "error=1") {
		t.Fatalf("unexpected Location %q (want /login?error=1)", loc)
	}
}

func TestHealth(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
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

func TestHealth_HEAD_returnsOK(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	req, err := http.NewRequest(http.MethodHead, srv.URL+"/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HEAD status %d want 200", resp.StatusCode)
	}
}

func TestStatic_cssServed(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	resp, err := http.Get(srv.URL + "/static/css/app.css")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	if len(b) < 50 {
		t.Fatalf("expected non-trivial css, got %d bytes", len(b))
	}
}

func TestStatic_unknownFileReturns404(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	resp, err := http.Get(srv.URL + "/static/moana-missing-asset-test-xyz.css")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d want 404", resp.StatusCode)
	}
}

func TestDashboardWithPeriodQuery(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "period@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "period@integration.test", "pw")
	resp, err := client.Get(srv.URL + "/?period=12m")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("dashboard status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `class="dashboard-hero-title"`) {
		t.Fatalf("expected dashboard shell")
	}
}

func TestDashboardWithUnknownPeriodQuery(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "period-unknown@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "period-unknown@integration.test", "pw")
	resp, err := client.Get(srv.URL + "/?period=weekly")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("dashboard status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `class="dashboard-hero-title"`) {
		t.Fatalf("expected dashboard shell (unknown period should fall back like parseStatsPeriod)")
	}
}

func TestLoginAndOverview(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "user@integration.test", "correct-password", "user")

	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "user@integration.test", "correct-password")

	resp, err := client.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("overview status %d", resp.StatusCode)
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

func TestTransactionsPageOKForLoggedInUser(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "txpage@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "txpage@moana.test", "pw")
	resp, err := client.Get(srv.URL + "/transactions")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "New entry") {
		t.Fatalf("expected new transaction form")
	}
}

func TestCategoriesPageOKForLoggedInUser(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "catpage@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "catpage@moana.test", "pw")
	resp, err := client.Get(srv.URL + "/categories")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Your categories") {
		t.Fatalf("expected categories page shell")
	}
}

func TestHistoryPageOKForLoggedInUser(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "hist@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "hist@moana.test", "pw")
	resp, err := client.Get(srv.URL + "/history")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `class="history-page"`) {
		t.Fatalf("expected history page shell")
	}
}

func TestHistoryPage_withQueryParams(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "histq@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "histq@moana.test", "pw")
	u := srv.URL + "/history?kind=expense&q=coffee&sort=oldest&from=2026-01-01&to=2026-01-31"
	resp, err := client.Get(u)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, `class="history-page"`) {
		t.Fatalf("expected history page shell")
	}
	if !strings.Contains(s, `name="q"`) {
		t.Fatalf("expected search field preserved")
	}
}

func TestHistoryPage_invalidDateRangeShowsBanner(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "histbad@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "histbad@moana.test", "pw")
	u := srv.URL + "/history?from=not-a-date&to=2020-01-02"
	resp, err := client.Get(u)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Invalid date range.") {
		t.Fatalf("expected date validation banner, got: %s", s[:min(600, len(s))])
	}
	if !strings.Contains(s, `class="alert alert-error"`) {
		t.Fatal("expected alert-error class")
	}
}

func TestSettingsPageOKForLoggedInUser(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "plain@moana.test", "pw", "user")
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
