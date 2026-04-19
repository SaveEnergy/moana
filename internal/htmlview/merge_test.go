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

func TestMergeFuncMaps_singleInputIsCloned(t *testing.T) {
	t.Parallel()
	orig := template.FuncMap{"n": func() int { return 1 }}
	m := MergeFuncMaps(orig)
	m["n"] = func() int { return 99 }
	f, _ := orig["n"].(func() int)
	if f == nil || f() != 1 {
		t.Fatalf("mutating merge result must not change source FuncMap")
	}
}

func TestMergeFuncMaps_twoNilMapsNonNilEmpty(t *testing.T) {
	t.Parallel()
	m := MergeFuncMaps(nil, nil)
	if m == nil {
		t.Fatal("expected non-nil empty map")
	}
	if len(m) != 0 {
		t.Fatalf("len %d", len(m))
	}
}

func TestMergeFuncMaps_nilFirstSecondSuppliesKeys(t *testing.T) {
	t.Parallel()
	m := MergeFuncMaps(nil, template.FuncMap{"a": func() int { return 7 }})
	af, _ := m["a"].(func() int)
	if af == nil || af() != 7 {
		t.Fatalf("got %T", m["a"])
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
