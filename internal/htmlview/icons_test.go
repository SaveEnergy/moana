package htmlview

import (
	"strings"
	"testing"
)

func TestIconFuncs_hasExpectedKeys(t *testing.T) {
	t.Parallel()
	m := IconFuncs()
	if len(m) != 3 {
		t.Fatalf("len %d", len(m))
	}
	for _, k := range []string{"moanaIcon", "moanaIconNav", "moanaIconPlain"} {
		if m[k] == nil {
			t.Fatalf("missing %q", k)
		}
	}
}

func TestMoanaIconPlain_knownID(t *testing.T) {
	t.Parallel()
	s := string(MoanaIconPlain("heart"))
	if !strings.Contains(s, "<svg") || !strings.Contains(s, `class="moana-icon"`) {
		t.Fatalf("got %q", s)
	}
}

func TestMoanaIconPlain_unknownIDUsesPaletteFallback(t *testing.T) {
	t.Parallel()
	s := string(MoanaIconPlain("not-a-registered-icon-id-xyz"))
	if len(s) < 30 || !strings.Contains(s, "<svg") {
		t.Fatalf("expected non-empty fallback SVG, got %q", s)
	}
}

func TestMoanaIconGrid_includesModifierClass(t *testing.T) {
	t.Parallel()
	s := string(MoanaIconGrid("heart"))
	if !strings.Contains(s, "moana-icon--grid") {
		t.Fatalf("got %q", s)
	}
}

func TestMoanaIconNav_includesModifierClass(t *testing.T) {
	t.Parallel()
	s := string(MoanaIconNav("heart"))
	if !strings.Contains(s, "moana-icon--nav") {
		t.Fatalf("got %q", s)
	}
}
