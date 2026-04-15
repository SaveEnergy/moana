package render

import (
	"bytes"
	"html/template"
	"net/http"

	"moana/internal/httperr"
)

// Shell executes the named page template (e.g. dashboard.html) into the layout body.
func (e *Engine) Shell(w http.ResponseWriter, contentTemplate string, data any, ld LayoutData, repoURL string) {
	var buf bytes.Buffer
	if err := e.Templates.ExecuteTemplate(&buf, contentTemplate, data); err != nil {
		httperr.Internal(w, nil, err)
		return
	}
	ld.RepoURL = repoURL
	ld.Body = template.HTML(buf.String())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := e.Templates.ExecuteTemplate(w, "layout.html", ld); err != nil {
		httperr.Internal(w, nil, err)
	}
}
