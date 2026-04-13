package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"moana/internal/store"
)

// LayoutData is the outer shell for authenticated pages.
type LayoutData struct {
	Title string
	User  *store.User
	Body  template.HTML
	Year  int
	Active string
	// MainClass is the CSS class for <main> (default: layer-stack in layout).
	MainClass string
}

func (a *App) renderShell(w http.ResponseWriter, innerName string, data any, ld LayoutData) {
	var buf bytes.Buffer
	if err := a.Templates.ExecuteTemplate(&buf, innerName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ld.Body = template.HTML(buf.String())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.Templates.ExecuteTemplate(w, "layout.html", ld); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
