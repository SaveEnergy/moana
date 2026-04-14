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
