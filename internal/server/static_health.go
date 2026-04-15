package server

import "net/http"

func registerStaticAndHealth(mux *http.ServeMux) {
	registerStatic(mux)
	registerHealth(mux)
}
