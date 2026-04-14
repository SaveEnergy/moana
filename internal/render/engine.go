package render

import (
	"bytes"
	"html/template"
	"net/http"

	"moana/internal/store"
)

// LayoutData is the outer shell for authenticated pages (layout.html).
type LayoutData struct {
	Title  string
	User   *store.User
	Body   template.HTML
	Year   int
	Active string
	// MainClass is an optional extra class on <main> (e.g. settings-shell).
	MainClass string
	// RepoURL is set by [Engine.Shell] from config (public GitHub / source link).
	RepoURL string
}

// Engine holds parsed templates and performs layout / simple renders.
type Engine struct {
	Templates *template.Template
}

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

// Simple executes a standalone template (e.g. login.html) without the app shell.
func (e *Engine) Simple(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := e.Templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
