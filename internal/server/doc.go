// Package server wires the HTTP mux: static assets with cache headers (static.go, static_test.go), health (health.go), request logging (logging.go, logging_test.go),
// bounded request bodies ([WithMaxRequestBodyBytes], maxbody.go, maxbody_test.go), per-request deadlines ([WithRequestTimeout], timeout_test.go), production [NewHTTPServer] (read/header/write/idle timeouts, bounded request headers; http_server.go, http_server_test.go), and application routes from [moana/internal/handlers] (router.go, router_test.go).
// [NewRouterWithRouterOptions] documents the outer→inner middleware chain (logging, optional timeout, max body, then mux).
// [router_test.go] includes [TestServeMux_GET_rootExactMatch] documenting why the dashboard uses `GET /{$}`.
// [ListenAndServeGracefully] runs the production server with graceful shutdown on SIGINT/SIGTERM (graceful.go, graceful_test.go for shutdown wiring).
package server
