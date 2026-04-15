package render

import (
	"net/http"

	"moana/internal/httperr"
)

// Simple executes a standalone template (e.g. login.html) without the app shell.
func (e *Engine) Simple(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := e.Templates.ExecuteTemplate(w, name, data); err != nil {
		httperr.Internal(w, nil, err)
	}
}
