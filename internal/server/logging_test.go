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

type flushCounter struct {
	*httptest.ResponseRecorder
	n int
}

func (c *flushCounter) Flush() {
	c.n++
	c.ResponseRecorder.Flush()
}

func TestStatusWriter_flushDelegatesToUnderlying(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	cf := &flushCounter{ResponseRecorder: rec}
	sw := &statusWriter{ResponseWriter: cf}
	sw.Flush()
	if cf.n != 1 {
		t.Fatalf("underlying Flush calls %d want 1", cf.n)
	}
}

type noFlusher struct{}

func (noFlusher) Header() http.Header       { return http.Header{} }
func (noFlusher) Write([]byte) (int, error) { return 0, nil }
func (noFlusher) WriteHeader(int)           {}

func TestStatusWriter_flushNoPanicWithoutFlusher(t *testing.T) {
	t.Parallel()
	sw := &statusWriter{ResponseWriter: noFlusher{}}
	sw.Flush()
}

func TestStatusWriter_unwrapReturnsUnderlying(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec}
	uw, ok := any(sw).(interface{ Unwrap() http.ResponseWriter })
	if !ok {
		t.Fatal("expected Unwrap method on statusWriter")
	}
	if uw.Unwrap() != rec {
		t.Fatal("Unwrap expected to return the wrapped ResponseWriter")
	}
}

func TestStatusWriter_pushReturnsErrNotSupportedWithoutPusher(t *testing.T) {
	t.Parallel()
	sw := &statusWriter{ResponseWriter: noFlusher{}}
	if err := sw.Push("/x", nil); err != http.ErrNotSupported {
		t.Fatalf("Push err %v want ErrNotSupported", err)
	}
}

type mockPusher struct {
	*httptest.ResponseRecorder
	lastTarget string
	lastOpts   *http.PushOptions
}

func (m *mockPusher) Push(target string, opts *http.PushOptions) error {
	m.lastTarget = target
	m.lastOpts = opts
	return nil
}

func TestStatusWriter_pushDelegatesToUnderlying(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	mp := &mockPusher{ResponseRecorder: rec}
	sw := &statusWriter{ResponseWriter: mp}
	opts := &http.PushOptions{Method: http.MethodGet}
	if err := sw.Push("/x.js", opts); err != nil {
		t.Fatalf("Push: %v", err)
	}
	if mp.lastTarget != "/x.js" {
		t.Fatalf("target %q", mp.lastTarget)
	}
	if mp.lastOpts != opts {
		t.Fatal("opts not forwarded")
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
