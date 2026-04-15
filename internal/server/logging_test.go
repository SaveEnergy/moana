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
