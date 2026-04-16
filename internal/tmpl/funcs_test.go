package tmpl

import (
	"net/url"
	"testing"
)

func TestTemplateFuncMap_urlQuery(t *testing.T) {
	t.Parallel()
	m := TemplateFuncMap()
	fn, ok := m["urlQuery"].(func(string) string)
	if !ok {
		t.Fatalf("urlQuery type got %T", m["urlQuery"])
	}
	want := url.QueryEscape("a b")
	if got := fn("a b"); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestTemplateFuncMap_mergesHtmlviewIconsAndHandlerFuncs(t *testing.T) {
	t.Parallel()
	m := TemplateFuncMap()
	for _, key := range []string{
		"formatEUR",
		"moanaIconPlain",
		"categoryPickerAccent",
		"urlQuery",
	} {
		if _, ok := m[key]; !ok {
			t.Fatalf("missing template func %q", key)
		}
	}
}
