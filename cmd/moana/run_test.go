package main

import "testing"

func TestRun_unknownSubcommand(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "nope"}); g != 1 {
		t.Fatalf("exit %d want 1", g)
	}
}

func TestRun_userWithoutSubcommand(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "user"}); g != 2 {
		t.Fatalf("exit %d want 2", g)
	}
}
