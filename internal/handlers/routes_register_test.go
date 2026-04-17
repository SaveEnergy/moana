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

// TestRegisterRoutes_notificationsGET verifies GET /notifications is registered (auth → redirect when anonymous).
func TestRegisterRoutes_notificationsGET(t *testing.T) {
	t.Parallel()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	a, err := moanaapp.New(testutil.DefaultTestConfig(), st)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, a)
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
