package httperr

import (
	"log/slog"
	"net/http"
)

// InternalMessage is the public response body for unexpected server failures (no internal detail).
const InternalMessage = "Something went wrong."

// Internal logs err and writes [InternalMessage] with status 500. If r is non-nil, method and path are logged.
func Internal(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if r != nil {
		slog.Error("http", "err", err, "method", r.Method, "path", r.URL.Path)
	} else {
		slog.Error("http", "err", err)
	}
	http.Error(w, InternalMessage, http.StatusInternalServerError)
}
