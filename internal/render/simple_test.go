package render

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"moana/internal/httperr"
)

func TestEngineSimple(t *testing.T) {
	t.Parallel()
	tmpl, err := template.New("page.html").Parse(`<!doctype html><p>{{.Title}}</p>`)
	if err != nil {
		t.Fatal(err)
	}
	e := &Engine{Templates: tmpl}
	rec := httptest.NewRecorder()
	e.Simple(rec, "page.html", struct{ Title string }{Title: "ok"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("Content-Type %q", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<p>ok</p>") {
		t.Fatalf("body %q", body)
	}
}

func TestSimple_execErrorDoesNotWritePartialHTML(t *testing.T) {
	t.Parallel()
	tmpl := template.Must(template.New("bad.html").Parse(`<!doctype html><p>{{.Missing}}</p>`))
	e := &Engine{Templates: tmpl}
	rec := httptest.NewRecorder()
	e.Simple(rec, "bad.html", struct{}{})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status %d want 500", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<p>") || strings.Contains(body, "<!doctype") {
		t.Fatalf("unexpected partial HTML: %q", body)
	}
	if !strings.Contains(body, httperr.InternalMessage) {
		t.Fatalf("expected internal message, got %q", body)
	}
}
