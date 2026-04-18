package server

import (
	"net/http"
	"testing"
	"time"
)

func TestNewHTTPServer_timeouts(t *testing.T) {
	t.Parallel()
	rt := 30 * time.Second
	srv := NewHTTPServer(":8080", rt, http.NotFoundHandler())
	if srv.Addr != ":8080" {
		t.Fatalf("Addr %q", srv.Addr)
	}
	if srv.ReadHeaderTimeout != readHeaderTimeout {
		t.Fatalf("ReadHeaderTimeout %v want %v", srv.ReadHeaderTimeout, readHeaderTimeout)
	}
	if srv.ReadTimeout != rt*2 {
		t.Fatalf("ReadTimeout %v want %v", srv.ReadTimeout, rt*2)
	}
	if srv.WriteTimeout != rt*2 {
		t.Fatalf("WriteTimeout %v want %v", srv.WriteTimeout, rt*2)
	}
	if srv.Handler == nil {
		t.Fatal("nil handler")
	}
}

func TestNewHTTPServer_zeroRequestTimeoutNoReadWriteDeadline(t *testing.T) {
	t.Parallel()
	srv := NewHTTPServer(":0", 0, http.NotFoundHandler())
	// Match [net/http.Server]: zero means no timeout on read/write body (distinct from ReadHeaderTimeout).
	if srv.ReadTimeout != 0 {
		t.Fatalf("ReadTimeout %v want 0", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 0 {
		t.Fatalf("WriteTimeout %v want 0", srv.WriteTimeout)
	}
	if srv.ReadHeaderTimeout != readHeaderTimeout {
		t.Fatalf("ReadHeaderTimeout still enforced: %v", srv.ReadHeaderTimeout)
	}
}
