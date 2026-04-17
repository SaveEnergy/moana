package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	moanaapp "moana/internal/app"
	"moana/internal/dbutil"
	"moana/internal/handlers"
	"moana/internal/testutil"
)

func newRegisterRoutesTestMux(t *testing.T) (*http.ServeMux, func()) {
	t.Helper()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	a, err := moanaapp.New(testutil.DefaultTestConfig(), st)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, a)
	return mux, func() { db.Close() }
}

func assertRedirectToLogin(t *testing.T, mux http.Handler, path string) {
	t.Helper()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("%s: status %d want %d", path, rec.Code, http.StatusSeeOther)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "/login") || !strings.Contains(loc, "error=1") {
		t.Fatalf("%s: Location %q want /login?...error=1...", path, loc)
	}
}

// TestRegisterRoutes_loginGET_ok verifies GET /login is registered and serves the login template (no auth).
func TestRegisterRoutes_loginGET_ok(t *testing.T) {
	t.Parallel()
	mux, cleanup := newRegisterRoutesTestMux(t)
	defer cleanup()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/login", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /login: status %d want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Sign in to your account") {
		t.Fatalf("expected login page HTML, body len %d", len(body))
	}
}

// TestRegisterRoutes_protectedGET_redirectsAnonymous verifies auth-wrapped GET routes redirect when there is no session.
func TestRegisterRoutes_protectedGET_redirectsAnonymous(t *testing.T) {
	t.Parallel()
	mux, cleanup := newRegisterRoutesTestMux(t)
	defer cleanup()
	paths := []string{
		"/",               // dashboard (GET /{$})
		"/transactions",
		"/transactions/42/edit", // path param (GET /transactions/{id}/edit)
		"/history",
		"/categories",
		"/settings",
		"/notifications",
	}
	for _, path := range paths {
		assertRedirectToLogin(t, mux, path)
	}
}

// TestRegisterRoutes_logoutPOST_redirectsToLogin verifies POST /logout is registered (no auth) and redirects to sign-in.
func TestRegisterRoutes_logoutPOST_redirectsToLogin(t *testing.T) {
	t.Parallel()
	mux, cleanup := newRegisterRoutesTestMux(t)
	defer cleanup()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/logout", nil))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("POST /logout: status %d want %d", rec.Code, http.StatusSeeOther)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "/login") {
		t.Fatalf("Location %q want /login", loc)
	}
}

// TestRegisterRoutes_unknownGET_404 verifies paths not matched by RegisterRoutes yield 404 (not the dashboard).
func TestRegisterRoutes_unknownGET_404(t *testing.T) {
	t.Parallel()
	mux, cleanup := newRegisterRoutesTestMux(t)
	defer cleanup()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/moana-route-register-test-unknown", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET unknown path: status %d want %d", rec.Code, http.StatusNotFound)
	}
}

// TestRegisterRoutes_protectedPOST_redirectsAnonymous verifies POST handlers wrapped in WithAuth redirect before touching the store.
func TestRegisterRoutes_protectedPOST_redirectsAnonymous(t *testing.T) {
	mux, cleanup := newRegisterRoutesTestMux(t)
	defer cleanup()
	cases := []struct {
		method, path string
	}{
		{http.MethodPost, "/transactions"},
		{http.MethodPost, "/transactions/42"},
		{http.MethodPost, "/categories"},
		{http.MethodPost, "/categories/update"},
		{http.MethodPost, "/categories/delete"},
		{http.MethodPost, "/settings/profile"},
		{http.MethodPost, "/settings/household"},
		{http.MethodPost, "/settings/household/members"},
		{http.MethodPost, "/settings/household/members/remove"},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != http.StatusSeeOther {
				t.Fatalf("status %d want %d", rec.Code, http.StatusSeeOther)
			}
			loc := rec.Header().Get("Location")
			if !strings.Contains(loc, "/login") || !strings.Contains(loc, "error=1") {
				t.Fatalf("Location %q want /login?...error=1...", loc)
			}
		})
	}
}
