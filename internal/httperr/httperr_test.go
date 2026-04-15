package httperr

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInternal_noLeak(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	Internal(w, r, errors.New("SECRET SQLITE BOILERPLATE"))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("code %d", w.Code)
	}
	body := w.Body.String()
	if strings.Contains(body, "SECRET") || strings.Contains(body, "SQLITE") {
		t.Fatalf("leaked detail: %q", body)
	}
	if !strings.Contains(body, InternalMessage) {
		t.Fatalf("body %q", body)
	}
}

func TestInternal_nilErrNoop(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	Internal(w, nil, nil)
	if w.Body.Len() != 0 {
		t.Fatal("expected no response body")
	}
}
