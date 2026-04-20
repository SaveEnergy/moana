package money

import (
	"math"
	"testing"
)

func TestAddCents(t *testing.T) {
	t.Parallel()
	if s, ok := AddCents(100, 200); !ok || s != 300 {
		t.Fatalf("got %d ok=%v", s, ok)
	}
	if s, ok := AddCents(0, math.MinInt64); !ok || s != math.MinInt64 {
		t.Fatalf("got %d ok=%v", s, ok)
	}
}

func TestAddCents_positiveOverflow(t *testing.T) {
	t.Parallel()
	_, ok := AddCents(math.MaxInt64, 1)
	if ok {
		t.Fatal("expected overflow")
	}
}

func TestAddCents_mergeMaxTwice(t *testing.T) {
	t.Parallel()
	_, ok := AddCents(math.MaxInt64, math.MaxInt64)
	if ok {
		t.Fatal("expected overflow")
	}
}
