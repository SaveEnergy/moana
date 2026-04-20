package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWithMaxRequestBodyBytes_parseFormTooLarge(t *testing.T) {
	t.Parallel()
	h := WithMaxRequestBodyBytes(100)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	body := strings.Repeat("a=", 60) // >100 bytes
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.ContentLength = int64(len(body))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code %d body %q", rec.Code, rec.Body.String())
	}
}

func TestWithMaxRequestBodyBytes_underLimitOK(t *testing.T) {
	t.Parallel()
	h := WithMaxRequestBodyBytes(1000)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("x=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d", rec.Code)
	}
}

func TestWithMaxRequestBodyBytes_zeroIsNoOp(t *testing.T) {
	t.Parallel()
	called := false
	h := WithMaxRequestBodyBytes(0)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rec, req)
	if !called {
		t.Fatal("handler not called")
	}
}
