package render

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"moana/internal/httperr"
)

func TestShell_writesFullPage(t *testing.T) {
	t.Parallel()
	tmpl := template.Must(template.New("").Parse(`
{{define "layout.html"}}<main>{{.Body}}</main>{{end}}
{{define "page.html"}}<p>hi</p>{{end}}
`))
	e := &Engine{Templates: tmpl}
	rec := httptest.NewRecorder()
	e.Shell(rec, "page.html", nil, LayoutData{Title: "t"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<main><p>hi</p></main>") {
		t.Fatalf("body %q", body)
	}
}

func TestShell_contentExecErrorDoesNotWriteHTML(t *testing.T) {
	t.Parallel()
	tmpl := template.Must(template.New("").Parse(`
{{define "layout.html"}}<main>{{.Body}}</main>{{end}}
{{define "page.html"}}<p>{{.Missing}}</p>{{end}}
`))
	e := &Engine{Templates: tmpl}
	rec := httptest.NewRecorder()
	e.Shell(rec, "page.html", struct{}{}, LayoutData{Title: "t"}, "")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status %d want 500", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<main>") || strings.Contains(body, "<p>") {
		t.Fatalf("unexpected partial HTML: %q", body)
	}
	if !strings.Contains(body, httperr.InternalMessage) {
		t.Fatalf("expected internal message, got %q", body)
	}
}

func TestShell_layoutExecErrorDoesNotWriteHTML(t *testing.T) {
	t.Parallel()
	tmpl := template.Must(template.New("").Parse(`
{{define "layout.html"}}<main>{{.Body}}{{.DoesNotExist}}</main>{{end}}
{{define "page.html"}}<p>x</p>{{end}}
`))
	e := &Engine{Templates: tmpl}
	rec := httptest.NewRecorder()
	e.Shell(rec, "page.html", nil, LayoutData{Title: "t"}, "")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status %d want 500", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<main>") {
		t.Fatalf("unexpected partial HTML: %q", body)
	}
	if !strings.Contains(body, httperr.InternalMessage) {
		t.Fatalf("expected internal message body, got %q", body)
	}
}
