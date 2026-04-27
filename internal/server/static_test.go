package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterStatic_faviconICORedirectsToStaticSVG(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	registerStatic(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/favicon.ico", nil))
	// [registerStatic] only registers /static/ — /favicon.ico is on [registerStaticAndHealth].
	if rec.Code != http.StatusNotFound {
		t.Fatalf("registerStatic only: GET /favicon.ico: status %d", rec.Code)
	}

	mux2 := http.NewServeMux()
	registerStaticAndHealth(mux2)
	rec2 := httptest.NewRecorder()
	mux2.ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, "/favicon.ico", nil))
	if rec2.Code != http.StatusFound {
		t.Fatalf("status %d want 302", rec2.Code)
	}
	if g, w := rec2.Header().Get("Location"), StaticURLPrefix+"favicon-mo.svg"; g != w {
		t.Fatalf("Location %q want %q", g, w)
	}
}

func TestRegisterStatic_servesCSSWithCacheControl(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	registerStatic(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, StaticURLPrefix+"css/app.css", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	cc := rec.Header().Get("Cache-Control")
	if !strings.Contains(cc, "max-age=") {
		t.Fatalf("Cache-Control %q", cc)
	}
	if cc != staticCacheControl {
		t.Fatalf("Cache-Control %q want %q", cc, staticCacheControl)
	}
}

func TestRegisterStatic_POSTReturns405(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	registerStatic(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, StaticURLPrefix+"css/app.css", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST %scss/app.css: status %d want 405", StaticURLPrefix, rec.Code)
	}
}

func TestStaticURLPrefix_matchesServeMuxStripPrefix(t *testing.T) {
	t.Parallel()
	const legacy = "GET /static/"
	if got := http.MethodGet + " " + StaticURLPrefix; got != legacy {
		t.Fatalf("registered pattern %q must stay %q for asset URLs", got, legacy)
	}
	if len(StaticURLPrefix) < 2 || StaticURLPrefix[len(StaticURLPrefix)-1] != '/' {
		t.Fatalf("StaticURLPrefix %q must end with / for StripPrefix + FileServer", StaticURLPrefix)
	}
}
