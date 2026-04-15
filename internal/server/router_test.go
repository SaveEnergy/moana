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
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
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
