package server

import "net/http"

// HealthPath is the GET liveness path for [registerHealth].
const HealthPath = "/health"

// healthOKBody is reused for every liveness response (avoids per-request []byte allocation).
var healthOKBody = []byte("ok")

func registerHealth(mux *http.ServeMux) {
	mux.HandleFunc(http.MethodGet+" "+HealthPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(healthOKBody)
	})
}
