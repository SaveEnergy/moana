package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHealth_returnsOKBody(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	registerHealth(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, HealthPath, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d", rec.Code)
	}
	if got := rec.Body.String(); got != "ok" {
		t.Fatalf("body %q want ok", got)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("Content-Type %q", ct)
	}
}

func TestHealthPath_matchesGo122ServeMuxPattern(t *testing.T) {
	t.Parallel()
	const legacyGETPattern = "GET /health"
	if got := http.MethodGet + " " + HealthPath; got != legacyGETPattern {
		t.Fatalf("registered pattern %q must stay %q for liveness probes", got, legacyGETPattern)
	}
}
