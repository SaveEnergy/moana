package server

import (
	"net/http"
	"time"
)

// readHeaderTimeout caps how long a client may take to send request headers (slowloris / stuck TCP).
// It is intentionally separate from handler/DB deadlines ([config.Config.RequestTimeout]).
const readHeaderTimeout = 10 * time.Second

// NewHTTPServer builds a production [http.Server] with connection-level timeouts.
// Read/write timeouts are derived from requestTimeout (request context deadline) but headers
// must arrive within readHeaderTimeout.
func NewHTTPServer(addr string, requestTimeout time.Duration, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       requestTimeout * 2,
		WriteTimeout:      requestTimeout * 2,
	}
}
