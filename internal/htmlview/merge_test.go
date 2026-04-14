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
