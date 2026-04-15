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
	if got := Internal(" /trim "); got != "/trim" {
		t.Fatalf("got %q", got)
	}
	// TrimSpace strips leading tab; path is still internal.
	if got := Internal("\t/foo"); got != "/foo" {
		t.Fatalf("got %q", got)
	}
	if got := Internal("\t//evil"); got != Default {
		t.Fatalf("got %q", got)
	}
	if got := Internal("/history?next=%2F"); got != "/history?next=%2F" {
		t.Fatalf("got %q", got)
	}
}
