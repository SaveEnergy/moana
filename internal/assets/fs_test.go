package assets

import (
	"io/fs"
	"testing"
)

func TestTemplateFS_hasHTMLTemplates(t *testing.T) {
	t.Parallel()
	fsys := TemplateFS()
	entries, err := fs.ReadDir(fsys, "templates")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected templates/*.html")
	}
	if _, err := fs.ReadFile(fsys, "templates/dashboard.html"); err != nil {
		t.Fatal(err)
	}
}

func TestStaticFS_servesAppCSS(t *testing.T) {
	t.Parallel()
	fsys, err := StaticFS()
	if err != nil {
		t.Fatal(err)
	}
	b, err := fs.ReadFile(fsys, "css/app.css")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) < 200 {
		t.Fatalf("expected non-trivial css, got %d bytes", len(b))
	}
}
