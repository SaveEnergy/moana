package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCategoryIDFromForm(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want int64
	}{
		{"", 0},
		{"abc", 0},
		{"0", 0},
		{"-3", 0},
		{"42", 42},
		{"  42  ", 42},
	}
	for _, tc := range cases {
		if got := categoryIDFromForm(tc.in); got != tc.want {
			t.Fatalf("categoryIDFromForm(%q) = %d want %d", tc.in, got, tc.want)
		}
	}
}

func TestPositiveInt64String(t *testing.T) {
	t.Parallel()
	if id, ok := positiveInt64String("  7  "); !ok || id != 7 {
		t.Fatalf("got %d ok=%v", id, ok)
	}
	if _, ok := positiveInt64String("  "); ok {
		t.Fatal("whitespace-only")
	}
}

func TestPositiveInt64String_maxInt64AndOverflow(t *testing.T) {
	t.Parallel()
	const maxI64 = int64(1<<63 - 1)
	tests := []struct {
		in     string
		wantID int64
		wantOK bool
	}{
		{"9223372036854775807", maxI64, true},
		{"9223372036854775808", 0, false}, // does not fit int64
		{"001", 1, true},
	}
	for _, tc := range tests {
		got, ok := positiveInt64String(tc.in)
		if ok != tc.wantOK || got != tc.wantID {
			t.Fatalf("positiveInt64String(%q) = (%d, %v), want (%d, %v)", tc.in, got, ok, tc.wantID, tc.wantOK)
		}
	}
}

func TestFormPositiveInt64(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("POST", CategoriesDeletePath, strings.NewReader("id=42&empty=&bad=0&neg=-1&space=+19+"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := r.ParseForm(); err != nil {
		t.Fatal(err)
	}
	if id, ok := formPositiveInt64(r, "id"); !ok || id != 42 {
		t.Fatalf("id: ok=%v id=%d", ok, id)
	}
	if id, ok := formPositiveInt64(r, "space"); !ok || id != 19 {
		t.Fatalf("trimmed id: ok=%v id=%d", ok, id)
	}
	for _, name := range []string{"empty", "bad", "neg", "missing"} {
		if _, ok := formPositiveInt64(r, name); ok {
			t.Fatalf("expected false for %q", name)
		}
	}
}
