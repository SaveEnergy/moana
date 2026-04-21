package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"moana/internal/config"
	"moana/internal/handlers"
)

// TestServeMux_GET_rootExactMatch documents stdlib Go 1.22+ behavior: [http.ServeMux] pattern "GET /"
// matches every path (prefix of "/"), while "GET /{$}" matches only "/". The app dashboard uses
// the latter so unknown paths are not served as the dashboard.
func TestServeMux_GET_rootExactMatch(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodGet+" "+handlers.DashboardRootPattern, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	t.Run("root", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Get(srv.URL + "/")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("got %d want 200", resp.StatusCode)
		}
	})
	t.Run("not_root", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Get(srv.URL + "/not-the-root")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("got %d want 404", resp.StatusCode)
		}
	})
}

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

func TestNewRouterWithRouterOptions_healthOKWithTimeoutAndMaxBodyStack(t *testing.T) {
	t.Parallel()
	// Minimal App: only Config is needed for GET /health (no store/render/session).
	cfg := &config.Config{RequestTimeout: 30 * time.Second}
	app := &handlers.App{Config: cfg}
	opts := &RouterOptions{DisableRequestLogging: true, MaxRequestBodyBytes: 4096}
	h := NewRouterWithRouterOptions(opts, app)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, HealthPath, nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d body %q", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body %q want ok", rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("Content-Type %q", ct)
	}
}

// TestNewRouterWithRouterOptions_healthNonGETReturns405 verifies [registerHealth] is GET-only on the full router stack.
func TestNewRouterWithRouterOptions_healthNonGETReturns405(t *testing.T) {
	cfg := &config.Config{RequestTimeout: 30 * time.Second}
	app := &handlers.App{Config: cfg}
	opts := &RouterOptions{DisableRequestLogging: true}
	h := NewRouterWithRouterOptions(opts, app)
	// HEAD is not listed: Go's ServeMux may route HEAD to GET handlers (implicit HEAD).
	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(method, HealthPath, nil))
			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s: status %d want 405", method, HealthPath, rec.Code)
			}
		})
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
