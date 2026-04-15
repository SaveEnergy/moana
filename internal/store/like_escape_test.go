package store

import "testing"

func TestEscapeSQLLikePattern(t *testing.T) {
	t.Parallel()
	if got := escapeSQLLikePattern("%"); got != "!%" {
		t.Fatalf("got %q want !%%", got)
	}
	if got := escapeSQLLikePattern("_"); got != "!_" {
		t.Fatalf("got %q", got)
	}
	if got := escapeSQLLikePattern("!"); got != "!!" {
		t.Fatalf("got %q", got)
	}
	if got := escapeSQLLikePattern("a!%_"); got != "a!!!%!_" {
		t.Fatalf("got %q", got)
	}
}
