package assets

import (
	"bytes"
	"io/fs"
	"testing"
)

func TestStaticFS_includesFaviconSVG(t *testing.T) {
	t.Parallel()
	sfs, err := StaticFS()
	if err != nil {
		t.Fatal(err)
	}
	b, err := fs.ReadFile(sfs, "favicon-mo.svg")
	if err != nil {
		t.Fatalf("embedded static must include favicon-mo.svg (served as /static/favicon-mo.svg): %v", err)
	}
	if len(b) < 50 {
		t.Fatalf("favicon-mo.svg too small: %d bytes", len(b))
	}
	if !bytes.Contains(b, []byte("<svg")) {
		t.Fatal("favicon-mo.svg must be SVG")
	}
}
