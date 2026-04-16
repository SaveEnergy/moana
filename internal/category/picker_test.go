package category_test

import (
	"testing"

	"moana/internal/category"
	"moana/internal/store"
)

func TestPickerAccent_validHexReturnsNormalized(t *testing.T) {
	t.Parallel()
	got := category.PickerAccent(store.Category{Color: "  #AaBbcc  ", ID: 99}, 0)
	if got != "#aabbcc" {
		t.Fatalf("got %q", got)
	}
}

func TestPickerAccent_invalidColorStableByCategoryID(t *testing.T) {
	t.Parallel()
	a := category.PickerAccent(store.Category{Color: "not-a-color", ID: 42}, 0)
	b := category.PickerAccent(store.Category{Color: "not-a-color", ID: 42}, 999)
	if a != b {
		t.Fatalf("same id should pick same accent: %q vs %q", a, b)
	}
}

func TestPickerAccent_invalidColorIDZeroUsesIndex(t *testing.T) {
	t.Parallel()
	x := category.PickerAccent(store.Category{Color: "bad", ID: 0}, 0)
	y := category.PickerAccent(store.Category{Color: "bad", ID: 0}, 1)
	if x == y {
		t.Fatal("expected different accents for different index when id is 0")
	}
}

func TestPickerIcon_keywordFood(t *testing.T) {
	t.Parallel()
	if got := category.PickerIcon("  Grocery day  ", 0); got != "utensils" {
		t.Fatalf("got %q", got)
	}
}

func TestPickerIcon_keywordTransport(t *testing.T) {
	t.Parallel()
	if got := category.PickerIcon("Fuel & parking", 0); got != "car" {
		t.Fatalf("got %q", got)
	}
}

func TestPickerIcon_defaultUsesIndex(t *testing.T) {
	t.Parallel()
	a := category.PickerIcon("qqqq", 0)
	b := category.PickerIcon("qqqq", 1)
	if a == b {
		t.Fatalf("fallback icons should vary by index: %q vs %q", a, b)
	}
}

func TestPickerHint_keywordAndDefault(t *testing.T) {
	t.Parallel()
	if got := category.PickerHint("  food  "); got != "Dining & groceries" {
		t.Fatalf("got %q", got)
	}
	if got := category.PickerHint("something utterly random"); got != "Activity in this category" {
		t.Fatalf("got %q", got)
	}
}
