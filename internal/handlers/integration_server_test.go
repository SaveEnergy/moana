package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"moana/internal/dashboard"
	"moana/internal/handlers"
	"moana/internal/historyview"
	"moana/internal/server"
	"moana/internal/testutil"
	"moana/internal/txform"
)

func TestLoginPageOK(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	resp, err := http.Get(srv.URL + handlers.LoginPath)
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
	if !strings.Contains(s, `class="login-oauth"`) || !strings.Contains(s, "OAuth sign-in (not available") {
		t.Fatalf("expected OAuth stub group with aria-label")
	}
}

func TestUnauthenticatedProtectedRoutesRedirectToLogin(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	for _, path := range []string{
		handlers.DashboardPath,
		handlers.HistoryPath,
		handlers.TransactionsPath,
		handlers.CategoriesPath,
		handlers.SettingsPath,
		handlers.NotificationsPath,
	} {
		resp, err := client.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusSeeOther {
			t.Fatalf("%s: status %d want %d", path, resp.StatusCode, http.StatusSeeOther)
		}
		loc := resp.Header.Get("Location")
		if !strings.Contains(loc, handlers.LoginPath) || !strings.Contains(loc, handlers.LoginErrorQueryParam+"=1") {
			t.Fatalf("%s: unexpected Location %q (want login+%s=1)", path, loc, handlers.LoginErrorQueryParam)
		}
	}
}

func TestHealth(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	resp, err := http.Get(srv.URL + server.HealthPath)
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
	req, err := http.NewRequest(http.MethodHead, srv.URL+server.HealthPath, nil)
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
	resp, err := http.Get(srv.URL + server.StaticURLPrefix + "css/app.css")
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
	resp, err := http.Get(srv.URL + server.StaticURLPrefix + "moana-missing-asset-test-xyz.css")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d want 404", resp.StatusCode)
	}
}

func TestUnknownAppRouteReturns404WhenUnauthenticated(t *testing.T) {
	t.Parallel()
	_, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	resp, err := http.Get(srv.URL + "/moana-no-route-anon-xyz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d want 404", resp.StatusCode)
	}
}

func TestUnknownAppRouteReturns404WhenAuthenticated(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "route404@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "route404@moana.test", "pw")
	resp, err := client.Get(srv.URL + "/moana-no-such-route-xyz")
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
	resp, err := client.Get(srv.URL + handlers.DashboardPath + "?" + dashboard.PeriodQueryParam + "=" + dashboard.StatsPeriod12m)
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
	resp, err := client.Get(srv.URL + handlers.DashboardPath + "?" + dashboard.PeriodQueryParam + "=weekly")
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

func TestDashboard_outflowShowsExpenseAfterCreate(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "dash-out@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "dash-out@integration.test", "pw")
	day := testutil.UTCDateString(time.Now())
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"42.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"dash outflow"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("post status %d", resp.StatusCode)
	}
	resp, err = client.Get(srv.URL + handlers.DashboardPath)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("dashboard status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, `class="dashboard-outflow-row"`) {
		t.Fatalf("expected outflow row, got: %s", s[:min(900, len(s))])
	}
	if !strings.Contains(s, "Uncategorized") {
		t.Fatalf("expected uncategorized category label, got: %s", s[:min(900, len(s))])
	}
	if !strings.Contains(s, "€42.00") {
		t.Fatalf("expected formatted amount in outflow, got: %s", s[:min(900, len(s))])
	}
}

func TestLoginAndOverview(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "user@integration.test", "correct-password", "user")

	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "user@integration.test", "correct-password")

	resp, err := client.Get(srv.URL + handlers.DashboardPath)
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
	resp, err := client.Get(srv.URL + handlers.TransactionsPath)
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
	resp, err := client.Get(srv.URL + handlers.CategoriesPath)
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
	resp, err := client.Get(srv.URL + handlers.HistoryPath)
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
	u := srv.URL + handlers.HistoryPath + "?kind=expense&q=coffee&sort=oldest&from=2026-01-01&to=2026-01-31"
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
	u := srv.URL + handlers.HistoryPath + "?from=not-a-date&to=2020-01-02"
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
	if !strings.Contains(s, historyview.InvalidDateRangeMessage) {
		t.Fatalf("expected date validation banner, got: %s", s[:min(600, len(s))])
	}
	assertBodyHasErrorAlert(t, s)
}

func TestHistoryPage_partialDateRangeShowsBanner(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "hist-partial@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "hist-partial@moana.test", "pw")
	for _, path := range []string{
		handlers.HistoryPath + "?from=2026-01-01",
		handlers.HistoryPath + "?to=2026-01-31",
	} {
		resp, err := client.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		func() {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("GET %s: status %d", path, resp.StatusCode)
			}
			body, _ := io.ReadAll(resp.Body)
			s := string(body)
			if !strings.Contains(s, historyview.InvalidDateRangeMessage) {
				t.Fatalf("GET %s: expected date validation copy, got: %s", path, s[:min(600, len(s))])
			}
			assertBodyHasErrorAlert(t, s)
		}()
	}
}

func TestNotificationsPageOKForLoggedInUser(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "notif@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "notif@moana.test", "pw")
	resp, err := client.Get(srv.URL + handlers.NotificationsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Notifications") || !strings.Contains(s, "no notifications") {
		t.Fatalf("expected notifications page shell, got prefix %q", s[:min(500, len(s))])
	}
}

func TestNotificationsPage_listsStoredNotifications(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	uid := testutil.MustCreateUser(t, app, "notif-listed@moana.test", "pw", "user")
	ctx := context.Background()
	if _, err := app.Store.InsertNotification(ctx, uid, "Integration test ping"); err != nil {
		t.Fatal(err)
	}
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "notif-listed@moana.test", "pw")
	resp, err := client.Get(srv.URL + handlers.NotificationsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Integration test ping") {
		t.Fatalf("expected notification body in HTML, got prefix %q", s[:min(800, len(s))])
	}
}

func TestSettingsPageOKForLoggedInUser(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "plain@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "plain@moana.test", "pw")
	resp, err := client.Get(srv.URL + handlers.SettingsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "settings-page") {
		t.Fatalf("expected settings page")
	}
	if !strings.Contains(s, "Custom profile photos") || !strings.Contains(s, "settings-photo-initial") {
		t.Fatalf("expected profile photo hint and initial, got prefix %q", s[:min(600, len(s))])
	}
}
