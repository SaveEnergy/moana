package server

import (
	"net/http"
	"time"
)

// readHeaderTimeout caps how long a client may take to send request headers (slowloris / stuck TCP).
// It is intentionally separate from handler/DB deadlines ([config.Config.RequestTimeout]).
const readHeaderTimeout = 10 * time.Second

// idleTimeout caps how long a keep-alive connection may sit idle between requests. Set explicitly so
// behavior does not depend on [http.Server] falling back to [http.Server.ReadTimeout] when IdleTimeout
// is zero (see net/http Server docs).
const idleTimeout = 120 * time.Second

// maxHeaderBytes caps request header size (matches [http.DefaultMaxHeaderBytes]); set explicitly so
// production limits are visible in code review.
const maxHeaderBytes = 1 << 20

// NewHTTPServer builds a production [http.Server] with connection-level timeouts and bounded headers
// ([maxHeaderBytes], same as [http.DefaultMaxHeaderBytes]).
// When requestTimeout is positive, read/write deadlines are 2× that value (room beyond the handler
// context deadline). Headers must still arrive within readHeaderTimeout. Keep-alive idle uses
// idleTimeout. When requestTimeout is zero, read/write/idle are left unset (no per-connection body
// deadline); ReadHeaderTimeout is still enforced.
func NewHTTPServer(addr string, requestTimeout time.Duration, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	if requestTimeout > 0 {
		d := requestTimeout * 2
		srv.ReadTimeout = d
		srv.WriteTimeout = d
		srv.IdleTimeout = idleTimeout
	}
	return srv
}
