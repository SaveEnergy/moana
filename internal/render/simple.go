package render

import (
	"bytes"
	"net/http"

	"moana/internal/httperr"
)

// Simple executes a standalone template (e.g. login.html) without the app shell.
// Output is buffered so template errors cannot emit partial HTML.
func (e *Engine) Simple(w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer
	if err := e.Templates.ExecuteTemplate(&buf, name, data); err != nil {
		httperr.Internal(w, nil, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}
