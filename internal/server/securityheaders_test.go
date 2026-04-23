package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithSecurityHeaders_setsBaseline(t *testing.T) {
	t.Parallel()
	h := WithSecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatal(rec.Code)
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options: %q", got)
	}
	if got := rec.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("X-Frame-Options: %q", got)
	}
	if rec.Header().Get("Referrer-Policy") == "" {
		t.Fatal("Referrer-Policy")
	}
	if rec.Header().Get("Permissions-Policy") == "" {
		t.Fatal("Permissions-Policy")
	}
}
