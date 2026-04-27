package server

import "net/http"

func registerStaticAndHealth(mux *http.ServeMux) {
	// Browsers, crawlers, and some proxies request /favicon.ico at the host root. Templates
	// use /static/favicon-mo.svg; redirect so the tab icon resolves without only rel=icon.
	mux.Handle("GET /favicon.ico", http.RedirectHandler(StaticURLPrefix+"favicon-mo.svg", http.StatusFound))
	registerStatic(mux)
	registerHealth(mux)
}
