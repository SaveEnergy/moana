package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusWriter_recordsWriteHeader(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec, status: http.StatusOK}
	sw.WriteHeader(http.StatusTeapot)
	if sw.status != http.StatusTeapot {
		t.Fatalf("status %d", sw.status)
	}
	if rec.Code != http.StatusTeapot {
		t.Fatalf("recorder code %d", rec.Code)
	}
}

func TestStatusWriter_defaultOKBeforeWriteHeader(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec, status: http.StatusOK}
	_, _ = sw.Write([]byte("ok"))
	if sw.status != http.StatusOK {
		t.Fatalf("status %d", sw.status)
	}
}

func TestRequestLogging_delegatesToInner(t *testing.T) {
	t.Parallel()
	var saw bool
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		saw = true
		w.WriteHeader(http.StatusTeapot)
	})
	rec := httptest.NewRecorder()
	RequestLogging(inner).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/tea", nil))
	if !saw {
		t.Fatal("inner handler not invoked")
	}
	if rec.Code != http.StatusTeapot {
		t.Fatalf("code %d want %d", rec.Code, http.StatusTeapot)
	}
}
