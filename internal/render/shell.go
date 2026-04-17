package render

import (
	"bytes"
	"html/template"
	"net/http"

	"moana/internal/httperr"
)

// Shell executes the named page template (e.g. dashboard.html) into the layout body.
// The full layout output is buffered before writing so template errors cannot emit partial HTML.
func (e *Engine) Shell(w http.ResponseWriter, contentTemplate string, data any, ld LayoutData, repoURL string) {
	var bodyBuf bytes.Buffer
	bodyBuf.Grow(4096)
	if err := e.Templates.ExecuteTemplate(&bodyBuf, contentTemplate, data); err != nil {
		httperr.Internal(w, nil, err)
		return
	}
	ld.RepoURL = repoURL
	ld.Body = template.HTML(bodyBuf.String())
	var out bytes.Buffer
	out.Grow(16384)
	if err := e.Templates.ExecuteTemplate(&out, "layout.html", ld); err != nil {
		httperr.Internal(w, nil, err)
		return
	}
	writeHTML(w, out.Bytes())
}
