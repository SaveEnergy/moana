package category_test

import (
	"testing"

	"moana/internal/category"
	"moana/internal/store"
)

func TestNormalizeColor(t *testing.T) {
	t.Parallel()
	if got := category.NormalizeColor("  #AbCdEf  "); got != "#abcdef" {
		t.Fatalf("got %q", got)
	}
	if got := category.NormalizeColor("not-a-color"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestIsCustomColor(t *testing.T) {
	t.Parallel()
	if !category.IsCustomColor(store.Category{Color: "#ff00ff"}) {
		t.Fatal("non-preset hex should be custom")
	}
	if category.IsCustomColor(store.Category{Color: "#ea580c"}) {
		t.Fatal("preset orange should not be custom")
	}
}

func TestTemplateColorPresets_matchesPresetsVar(t *testing.T) {
	t.Parallel()
	got := category.TemplateColorPresets()
	if len(got) != len(category.Presets) {
		t.Fatalf("len %d vs Presets %d", len(got), len(category.Presets))
	}
	if len(got) == 0 || got[0].Hex != category.Presets[0].Hex {
		t.Fatalf("first preset mismatch")
	}
}

func TestHexOrDefault(t *testing.T) {
	t.Parallel()
	if got := category.HexOrDefault(store.Category{Color: "#AaBbcC"}); got != "#aabbcc" {
		t.Fatalf("got %q", got)
	}
	if got := category.HexOrDefault(store.Category{Color: ""}); got != "#818cf8" {
		t.Fatalf("default periwinkle: got %q", got)
	}
	if got := category.HexOrDefault(store.Category{Color: "not-a-color"}); got != "#818cf8" {
		t.Fatalf("invalid falls back: got %q", got)
	}
}

