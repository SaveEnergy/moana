package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"moana/internal/assets"
)

// staticCacheControl is applied to embedded /static/* responses. Paths are not
// content-hashed, so we avoid immutable long-term caching; repeat visits still
// benefit from a modest max-age.
const staticCacheControl = "public, max-age=86400"

func registerStatic(mux *http.ServeMux) {
	staticFS, err := assets.StaticFS()
	if err != nil {
		panic(fmt.Errorf("server: embedded static assets: %w", err))
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticFileHandler(staticFS)))
}

func staticFileHandler(root fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", staticCacheControl)
		fileServer.ServeHTTP(w, r)
	})
}
