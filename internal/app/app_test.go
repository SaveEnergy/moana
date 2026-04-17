package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	moanaapp "moana/internal/app"
	"moana/internal/dbutil"
	"moana/internal/testutil"
)

func TestNew_parsesTemplates(t *testing.T) {
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
	if a == nil || a.Render == nil || a.Render.Templates == nil {
		t.Fatal("expected app with render engine")
	}
}

func TestHTTPHandler_wiresRouter(t *testing.T) {
	t.Parallel()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	h, err := moanaapp.HTTPHandler(testutil.DefaultTestConfig(), st)
	if err != nil {
		t.Fatal(err)
	}
	if h == nil {
		t.Fatal("nil handler")
	}
}

func TestHTTPHandler_servesHealth(t *testing.T) {
	t.Parallel()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	h, err := moanaapp.HTTPHandler(testutil.DefaultTestConfig(), st)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body %q", rec.Body.String())
	}
}

func TestHTTPHandler_servesStaticCSS(t *testing.T) {
	t.Parallel()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	h, err := moanaapp.HTTPHandler(testutil.DefaultTestConfig(), st)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/css/app.css", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if n := rec.Body.Len(); n < 100 {
		t.Fatalf("expected non-trivial css, got %d bytes", n)
	}
}
