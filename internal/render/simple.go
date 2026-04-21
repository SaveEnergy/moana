package render

import (
	"bytes"
	"net/http"

	"moana/internal/httperr"
)

// Simple executes a standalone template (e.g. login.html) without the app shell.
// Output is buffered so template errors cannot emit partial HTML.
// r is passed to [httperr.Internal] on failure so logs include method and path (r may be nil in tests).
func (e *Engine) Simple(w http.ResponseWriter, r *http.Request, name string, data any) {
	var buf bytes.Buffer
	buf.Grow(4096)
	if err := e.Templates.ExecuteTemplate(&buf, name, data); err != nil {
		httperr.Internal(w, r, err)
		return
	}
	writeHTML(w, buf.Bytes())
}
