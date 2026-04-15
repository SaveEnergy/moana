package server

import (
	"fmt"
	"net/http"

	"moana/internal/assets"
)

func registerStatic(mux *http.ServeMux) {
	staticFS, err := assets.StaticFS()
	if err != nil {
		panic(fmt.Errorf("server: embedded static assets: %w", err))
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
}
