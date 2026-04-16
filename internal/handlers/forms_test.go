package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"
)

func TestRequireParseForm_success(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("a=b&c=d"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	if !requireParseForm(rec, req) {
		t.Fatal("expected true")
	}
	if req.FormValue("a") != "b" || req.FormValue("c") != "d" {
		t.Fatalf("form values a=%q c=%q", req.FormValue("a"), req.FormValue("c"))
	}
}

func TestRequireParseForm_bodyReadError(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/", iotest.ErrReader(errors.New("read fail")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	if requireParseForm(rec, req) {
		t.Fatal("expected false")
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "bad form") {
		t.Fatalf("body %q", rec.Body.String())
	}
}
