package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"moana/internal/config"
	"moana/internal/handlers"
)

func TestRequestTimeout_prefersRouterOptions(t *testing.T) {
	t.Parallel()
	app := &handlers.App{Config: &config.Config{RequestTimeout: 5 * time.Second}}
	opts := &RouterOptions{RequestTimeout: 10 * time.Second}
	if got := requestTimeout(opts, app); got != 10*time.Second {
		t.Fatalf("got %v", got)
	}
}

func TestRequestTimeout_usesAppConfigWhenOptsUnset(t *testing.T) {
	t.Parallel()
	app := &handlers.App{Config: &config.Config{RequestTimeout: 7 * time.Second}}
	if got := requestTimeout(nil, app); got != 7*time.Second {
		t.Fatalf("got %v", got)
	}
}

func TestRequestTimeout_zeroWhenNoConfig(t *testing.T) {
	t.Parallel()
	if got := requestTimeout(nil, &handlers.App{}); got != 0 {
		t.Fatalf("got %v", got)
	}
	if got := requestTimeout(&RouterOptions{}, nil); got != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestMaxRequestBodyBytes_resolution(t *testing.T) {
	t.Parallel()
	if got := maxRequestBodyBytes(nil); got != defaultMaxRequestBodyBytes {
		t.Fatalf("nil opts: got %d", got)
	}
	if got := maxRequestBodyBytes(&RouterOptions{}); got != defaultMaxRequestBodyBytes {
		t.Fatalf("zero field: got %d", got)
	}
	const custom = int64(4096)
	if got := maxRequestBodyBytes(&RouterOptions{MaxRequestBodyBytes: custom}); got != custom {
		t.Fatalf("override: got %d want %d", got, custom)
	}
	if got := maxRequestBodyBytes(&RouterOptions{MaxRequestBodyBytes: -1}); got != 0 {
		t.Fatalf("negative disables cap: got %d want 0", got)
	}
}

func TestRouterMiddlewareComposition_matchesNewRouterWithRouterOptions(t *testing.T) {
	t.Parallel()
	// Same stack as [NewRouterWithRouterOptions] minus [RequestLogging]: [wrapRouterMiddleware].
	var sawDeadline bool
	mux := http.NewServeMux()
	mux.HandleFunc("GET /probe", func(w http.ResponseWriter, r *http.Request) {
		_, sawDeadline = r.Context().Deadline()
		w.WriteHeader(http.StatusOK)
	})
	cfg := &config.Config{RequestTimeout: 30 * time.Second}
	app := &handlers.App{Config: cfg}
	opts := &RouterOptions{MaxRequestBodyBytes: 100}
	h := wrapRouterMiddleware(mux, opts, app)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d", rec.Code)
	}
	if !sawDeadline {
		t.Fatal("handler must see request context deadline (wrapRouterMiddleware must stay in sync with NewRouterWithRouterOptions)")
	}
}

func TestWrapRouterMiddleware_appliesTimeoutFromAppWhenOptsUnset(t *testing.T) {
	t.Parallel()
	var sawDeadline bool
	mux := http.NewServeMux()
	mux.HandleFunc("GET /probe", func(w http.ResponseWriter, r *http.Request) {
		_, sawDeadline = r.Context().Deadline()
		w.WriteHeader(http.StatusOK)
	})
	app := &handlers.App{Config: &config.Config{RequestTimeout: 5 * time.Second}}
	opts := &RouterOptions{MaxRequestBodyBytes: 512}
	h := wrapRouterMiddleware(mux, opts, app)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))
	if !sawDeadline {
		t.Fatal("expected deadline from app.Config.RequestTimeout when RouterOptions.RequestTimeout is zero")
	}
}

func TestWrapRouterMiddleware_noDeadlineWhenTimeoutUnset(t *testing.T) {
	t.Parallel()
	var sawDeadline bool
	mux := http.NewServeMux()
	mux.HandleFunc("GET /probe", func(w http.ResponseWriter, r *http.Request) {
		_, sawDeadline = r.Context().Deadline()
		w.WriteHeader(http.StatusOK)
	})
	h := wrapRouterMiddleware(mux, &RouterOptions{}, &handlers.App{})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d", rec.Code)
	}
	if sawDeadline {
		t.Fatal("expected no request deadline when RequestTimeout is unset (WithRequestTimeout must be skipped)")
	}
}
