// Package server wires the HTTP mux: static assets, health check, request logging, and
// application routes from [moana/internal/handlers]. [ListenAndServeGracefully] runs the
// production server with graceful shutdown on SIGINT/SIGTERM.
package server
