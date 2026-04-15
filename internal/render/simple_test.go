package render

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	body := rec.Body.String()
	if !strings.Contains(body, "<p>ok</p>") {
		t.Fatalf("body %q", body)
	}
}
