package server

import (
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
