package htmlview

import (
	"html/template"
	"testing"
)

func TestMergeFuncMaps(t *testing.T) {
	t.Parallel()
	a := template.FuncMap{"x": func() int { return 1 }}
	b := template.FuncMap{"y": func() int { return 2 }}
	c := template.FuncMap{"x": func() int { return 3 }}
	m := MergeFuncMaps(a, b, c)
	if len(m) != 2 {
		t.Fatalf("len=%d", len(m))
	}
	fn := m["x"].(func() int)
	if fn() != 3 {
		t.Fatal("last x should win")
	}
}

func TestMergeFuncMaps_capacityMatchesKeys(t *testing.T) {
	t.Parallel()
	a := template.FuncMap{"a": func() {}, "b": func() {}}
	b := template.FuncMap{"c": func() {}}
	m := MergeFuncMaps(a, b)
	if len(m) != 3 {
		t.Fatalf("len=%d", len(m))
	}
}
