package htmlview

import (
	"html/template"
	"testing"
)

func TestMergeFuncMaps_laterMapOverrides(t *testing.T) {
	t.Parallel()
	m := MergeFuncMaps(
		template.FuncMap{"x": func() string { return "first" }},
		template.FuncMap{"x": func() string { return "second" }},
	)
	fn, ok := m["x"].(func() string)
	if !ok {
		t.Fatalf("merged x type %T", m["x"])
	}
	if got := fn(); got != "second" {
		t.Fatalf("got %q want second (rightmost map wins)", got)
	}
}

func TestMergeFuncMaps_emptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := MergeFuncMaps()
	if m == nil {
		t.Fatal("expected non-nil map")
	}
	if len(m) != 0 {
		t.Fatalf("len %d", len(m))
	}
}

func TestMergeFuncMaps_mergesDistinctKeys(t *testing.T) {
	t.Parallel()
	m := MergeFuncMaps(
		template.FuncMap{"a": func() int { return 1 }},
		template.FuncMap{"b": func() int { return 2 }},
	)
	af, _ := m["a"].(func() int)
	bf, _ := m["b"].(func() int)
	if af == nil || bf == nil {
		t.Fatalf("missing funcs: a=%T b=%T", m["a"], m["b"])
	}
	if af() != 1 || bf() != 2 {
		t.Fatalf("a=%d b=%d", af(), bf())
	}
}

func BenchmarkMergeFuncMaps_threeMaps(b *testing.B) {
	a := template.FuncMap{"a": func() int { return 1 }}
	x := template.FuncMap{"b": func() int { return 2 }}
	z := template.FuncMap{"c": func() int { return 3 }}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = MergeFuncMaps(a, x, z)
	}
}
