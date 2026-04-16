package icons

import (
	"strings"
	"testing"
)

func TestValidID(t *testing.T) {
	t.Parallel()
	if ValidID("") {
		t.Fatal("empty id")
	}
	if !ValidID("heart") {
		t.Fatal("known id")
	}
	if ValidID("not-a-real-lucide-id-xyz") {
		t.Fatal("unknown id")
	}
}

func TestSVG_rendersKnownIcon(t *testing.T) {
	t.Parallel()
	h := SVG("heart", "icon-test")
	if !strings.Contains(string(h), `<svg`) || !strings.Contains(string(h), `class="icon-test"`) {
		t.Fatalf("unexpected: %s", h)
	}
}

func TestSVG_unknownFallsBackToPalette(t *testing.T) {
	t.Parallel()
	h := SVG("totally-unknown-icon-id", "x")
	if !strings.Contains(string(h), `<svg`) {
		t.Fatal("expected fallback svg")
	}
}

func TestSVG_escapesClassAttribute(t *testing.T) {
	t.Parallel()
	// Class is interpolated into HTML; must not allow breaking out of the attribute.
	h := SVG("heart", `"><img src=x`)
	s := string(h)
	if strings.Contains(s, `"><img`) {
		t.Fatalf("class not escaped: %s", s)
	}
	if !strings.Contains(s, `class="`) {
		t.Fatalf("missing class attr: %s", s)
	}
}

func TestInner(t *testing.T) {
	t.Parallel()
	if Inner("heart") == "" {
		t.Fatal("expected path data")
	}
	if Inner("not-a-real-lucide-id-xyz") != "" {
		t.Fatal("unknown should be empty")
	}
	if Inner("") != "" {
		t.Fatal("empty id should be empty inner")
	}
}
