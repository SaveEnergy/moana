package render

import (
	"bytes"
	"html/template"
	"net/http"
)

// Shell executes the named page template (e.g. dashboard.html) into the layout body.
func (e *Engine) Shell(w http.ResponseWriter, contentTemplate string, data any, ld LayoutData, repoURL string) {
	var buf bytes.Buffer
	if err := e.Templates.ExecuteTemplate(&buf, contentTemplate, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ld.RepoURL = repoURL
	ld.Body = template.HTML(buf.String())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := e.Templates.ExecuteTemplate(w, "layout.html", ld); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
