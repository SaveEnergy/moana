package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPathPositiveInt64(t *testing.T) {
	t.Parallel()
	cases := []struct {
		val  string
		want int64
		ok   bool
	}{
		{"42", 42, true},
		{"1", 1, true},
		{"0", 0, false},
		{"-3", 0, false},
		{"", 0, false},
		{"nope", 0, false},
	}
	for _, tc := range cases {
		r := httptest.NewRequest(http.MethodGet, "/transactions/1/edit", nil)
		r.SetPathValue("id", tc.val)
		got, ok := pathPositiveInt64(r, "id")
		if ok != tc.ok || got != tc.want {
			t.Fatalf("id=%q: got (%d,%v) want (%d,%v)", tc.val, got, ok, tc.want, tc.ok)
		}
	}
}
