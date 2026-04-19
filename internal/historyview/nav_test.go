package historyview

import "testing"

func TestPathWithQuery(t *testing.T) {
	t.Parallel()
	if got := pathWithQuery(""); got != RoutePath {
		t.Fatalf("empty encoded query: got %q want %q", got, RoutePath)
	}
	const q = "q=a&kind=expense"
	if got, want := pathWithQuery(q), RoutePath+"?"+q; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
