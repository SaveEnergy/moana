package dashboard

import (
	"net/url"
	"strings"
	"testing"
)

func TestPeriodQueryRelative_escapesAndUsesParamName(t *testing.T) {
	t.Parallel()
	got := PeriodQueryRelative(StatsPeriod30d)
	if !strings.HasPrefix(got, "/?"+PeriodQueryParam+"=") {
		t.Fatalf("got %q", got)
	}
	v, err := url.ParseQuery(strings.TrimPrefix(got, "/?"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Get(PeriodQueryParam) != StatsPeriod30d {
		t.Fatalf("query %+v", v)
	}
}

func TestPeriodQueryRelative_escapesUnsafeQueryValue(t *testing.T) {
	t.Parallel()
	got := PeriodQueryRelative("a&b=c")
	v, err := url.ParseQuery(strings.TrimPrefix(got, "/?"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Get(PeriodQueryParam) != "a&b=c" {
		t.Fatalf("round-trip got %q query=%q", v.Get(PeriodQueryParam), got)
	}
}
