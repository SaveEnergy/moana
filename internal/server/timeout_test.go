package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWithRequestTimeout_zeroIsPassthrough(t *testing.T) {
	t.Parallel()
	var ran bool
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ran = true
		if _, ok := r.Context().Deadline(); ok {
			t.Fatal("unexpected deadline")
		}
	})
	h := WithRequestTimeout(0)(inner)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !ran {
		t.Fatal("inner not called")
	}
}

func TestWithRequestTimeout_setsDeadline(t *testing.T) {
	t.Parallel()
	var hasDeadline bool
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasDeadline = r.Context().Deadline()
	})
	h := WithRequestTimeout(30 * time.Second)(inner)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !hasDeadline {
		t.Fatal("expected request context deadline")
	}
}

func TestWithRequestTimeout_cancelsAfterExpiry(t *testing.T) {
	t.Parallel()
	done := make(chan struct{})
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		close(done)
	})
	h := WithRequestTimeout(1*time.Millisecond)(inner)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not finish after context cancel")
	}
}
