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
	buf.Grow(4096)
	if err := e.Templates.ExecuteTemplate(&buf, name, data); err != nil {
		httperr.Internal(w, nil, err)
		return
	}
	writeHTML(w, buf.Bytes())
}
