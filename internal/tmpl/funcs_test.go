package tmpl

import (
	"net/url"
	"testing"

	"moana/internal/dashboard"
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
		"statsPeriod30d",
		"statsPeriod12m",
		"dashboardPeriodHref",
	} {
		if _, ok := m[key]; !ok {
			t.Fatalf("missing template func %q", key)
		}
	}
}

func TestTemplateFuncMap_dashboardPeriodHref(t *testing.T) {
	t.Parallel()
	m := TemplateFuncMap()
	fn, ok := m["dashboardPeriodHref"].(func(string) string)
	if !ok {
		t.Fatalf("dashboardPeriodHref type got %T", m["dashboardPeriodHref"])
	}
	s30, ok := m["statsPeriod30d"].(func() string)
	if !ok {
		t.Fatalf("statsPeriod30d type got %T", m["statsPeriod30d"])
	}
	want := dashboard.PeriodQueryRelative(dashboard.StatsPeriod30d)
	if got := fn(s30()); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
