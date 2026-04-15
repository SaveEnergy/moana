package handlers

import "testing"

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
	}
	for _, tc := range cases {
		if got := categoryIDFromForm(tc.in); got != tc.want {
			t.Fatalf("categoryIDFromForm(%q) = %d want %d", tc.in, got, tc.want)
		}
	}
}
