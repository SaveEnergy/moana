package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

// TestNewHTTPServer_ServeShutdown_returnsErrServerClosed documents that production wiring from
// [NewHTTPServer] supports graceful shutdown: [http.Server.Shutdown] stops [http.Server.Serve].
func TestNewHTTPServer_ServeShutdown_returnsErrServerClosed(t *testing.T) {
	t.Parallel()
	srv := NewHTTPServer("127.0.0.1:0", 30*time.Second, http.NotFoundHandler())
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Serve(ln) }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	recv := <-errCh
	if recv != http.ErrServerClosed {
		t.Fatalf("Serve returned %v want %v", recv, http.ErrServerClosed)
	}
}
