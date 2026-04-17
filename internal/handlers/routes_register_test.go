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

// TestRegisterRoutes_notificationsGET verifies GET /notifications is registered (auth → redirect when anonymous).
func TestRegisterRoutes_notificationsGET(t *testing.T) {
	t.Parallel()
	mux, cleanup := newRegisterRoutesTestMux(t)
	defer cleanup()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/notifications", nil))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("GET /notifications: status %d want %d", rec.Code, http.StatusSeeOther)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "/login") || !strings.Contains(loc, "error=1") {
		t.Fatalf("Location %q want /login?...error=1...", loc)
	}
}
