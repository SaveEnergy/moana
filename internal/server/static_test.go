package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
