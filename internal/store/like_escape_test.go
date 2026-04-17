package store

import "testing"

func TestEscapeSQLLikePattern_metachars(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want string
	}{
		{"plain", "plain"},
		{"a%b", "a!%b"},
		{"a_b", "a!_b"},
		{"!", "!!"},
		{"%_%", "!%!_!%"},
	}
	for _, tc := range cases {
		if got := escapeSQLLikePattern(tc.in); got != tc.want {
			t.Fatalf("escapeSQLLikePattern(%q) = %q want %q", tc.in, got, tc.want)
		}
	}
}
