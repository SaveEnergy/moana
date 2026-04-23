package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"moana/internal/config"
	"moana/internal/handlers"
)

func TestRollingMinuteLimiter_allow(t *testing.T) {
	t.Parallel()
	lim := newRollingMinuteLimiter(2, 100*time.Millisecond)
	if !lim.allow("a") {
		t.Fatal("first")
	}
	if !lim.allow("a") {
		t.Fatal("second")
	}
	if lim.allow("a") {
		t.Fatal("third should block")
	}
}

func TestWithPostAuthRateLimit_blocksLogin(t *testing.T) {
	cfg := &config.Config{
		RateLimitLoginPerMin: 2,
	}
	app := &handlers.App{Config: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+handlers.LoginPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(WithPostAuthRateLimit(app)(mux))
	t.Cleanup(srv.Close)
	client := srv.Client()
	for i := 0; i < 2; i++ {
		resp, err := client.Post(srv.URL+handlers.LoginPath, "application/x-www-form-urlencoded", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("req %d: %d", i, resp.StatusCode)
		}
	}
	resp, err := client.Post(srv.URL+handlers.LoginPath, "application/x-www-form-urlencoded", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("code %d want 429", resp.StatusCode)
	}
}

func TestWithPostAuthRateLimit_usesXForwardedForWhenTrusted(t *testing.T) {
	cfg := &config.Config{
		RateLimitLoginPerMin: 1,
		TrustForwardedAddr:   true,
	}
	app := &handlers.App{Config: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+handlers.LoginPath, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := WithPostAuthRateLimit(app)(mux)

	req1 := httptest.NewRequest(http.MethodPost, handlers.LoginPath, nil)
	req1.Header.Set("X-Forwarded-For", "10.0.0.1, 1.1.1.1")
	req1.RemoteAddr = "198.18.0.1:9"
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first: %d", rec1.Code)
	}

	req2 := httptest.NewRequest(http.MethodPost, handlers.LoginPath, nil)
	req2.Header.Set("X-Forwarded-For", "10.0.0.1")
	req2.RemoteAddr = "198.18.0.1:9"
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("same XFF client: %d want 429", rec2.Code)
	}

	req3 := httptest.NewRequest(http.MethodPost, handlers.LoginPath, nil)
	req3.Header.Set("X-Forwarded-For", "10.0.0.2")
	req3.RemoteAddr = "198.18.0.1:9"
	rec3 := httptest.NewRecorder()
	h.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		t.Fatalf("other client: %d", rec3.Code)
	}
}
