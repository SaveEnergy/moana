package render

import (
	"bytes"
	"html/template"
	"net/http"

	"moana/internal/httperr"
)

// Shell executes the named page template (e.g. dashboard.html) into the layout body.
// The full layout output is buffered before writing so template errors cannot emit partial HTML.
// r is passed to [httperr.Internal] on failure so logs include method and path (r may be nil in tests).
func (e *Engine) Shell(w http.ResponseWriter, r *http.Request, contentTemplate string, data any, ld LayoutData, repoURL string) {
	var bodyBuf bytes.Buffer
	bodyBuf.Grow(4096)
	if err := e.Templates.ExecuteTemplate(&bodyBuf, contentTemplate, data); err != nil {
		httperr.Internal(w, r, err)
		return
	}
	ld.RepoURL = repoURL
	ld.Body = template.HTML(bodyBuf.String())
	var out bytes.Buffer
	out.Grow(16384)
	if err := e.Templates.ExecuteTemplate(&out, "layout.html", ld); err != nil {
		httperr.Internal(w, r, err)
		return
	}
	writeHTML(w, out.Bytes())
}
