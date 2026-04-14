package safepath

import "testing"

func TestInternal(t *testing.T) {
	t.Parallel()
	if got := Internal(""); got != Default {
		t.Fatal(got)
	}
	if got := Internal("/foo/bar"); got != "/foo/bar" {
		t.Fatal(got)
	}
	if got := Internal("//evil.com"); got != Default {
		t.Fatal(got)
	}
	if got := Internal("https://x"); got != Default {
		t.Fatal(got)
	}
}
