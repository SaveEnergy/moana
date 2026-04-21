package app

import (
	"testing"
	"time"

	"moana/internal/config"
)

func TestRouterOptionsFromConfig_nilReturnsNil(t *testing.T) {
	t.Parallel()
	if got := routerOptionsFromConfig(nil); got != nil {
		t.Fatalf("got %#v want nil", got)
	}
}

func TestRouterOptionsFromConfig_mapsTimeoutAndMaxBody(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		RequestTimeout:      42 * time.Second,
		MaxRequestBodyBytes: 8192,
	}
	opts := routerOptionsFromConfig(cfg)
	if opts == nil {
		t.Fatal("nil opts")
	}
	if opts.RequestTimeout != cfg.RequestTimeout {
		t.Fatalf("RequestTimeout %v want %v", opts.RequestTimeout, cfg.RequestTimeout)
	}
	if opts.MaxRequestBodyBytes != cfg.MaxRequestBodyBytes {
		t.Fatalf("MaxRequestBodyBytes %d want %d", opts.MaxRequestBodyBytes, cfg.MaxRequestBodyBytes)
	}
}

func TestRouterOptionsFromConfig_zeroValuesStillSet(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{}
	opts := routerOptionsFromConfig(cfg)
	if opts == nil {
		t.Fatal("nil opts")
	}
	if opts.RequestTimeout != 0 || opts.MaxRequestBodyBytes != 0 {
		t.Fatalf("got timeout=%v maxBody=%d want zeros", opts.RequestTimeout, opts.MaxRequestBodyBytes)
	}
}
