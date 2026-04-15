// Package server wires the HTTP mux: static assets (static.go), health (health.go), request logging (logging.go, logging_test.go),
// per-request deadlines ([WithRequestTimeout], timeout_test.go), and application routes from [moana/internal/handlers] (router.go, router_test.go).
// [router_test.go] includes [TestServeMux_GET_rootExactMatch] documenting why the dashboard uses `GET /{$}`.
// [ListenAndServeGracefully] runs the production server with graceful shutdown on SIGINT/SIGTERM (graceful.go).
package server
