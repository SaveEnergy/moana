package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
