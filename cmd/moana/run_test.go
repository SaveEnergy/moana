package main

import "testing"

func TestRun_version(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "version"}); g != 0 {
		t.Fatalf("exit %d want 0", g)
	}
}

func TestRun_versionExtraArgs(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "version", "extra"}); g != 1 {
		t.Fatalf("exit %d want 1", g)
	}
}

func TestRun_dashVersion(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "-version"}); g != 0 {
		t.Fatalf("exit %d want 0", g)
	}
}

func TestRun_doubleDashVersion(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "--version"}); g != 0 {
		t.Fatalf("exit %d want 0", g)
	}
}

func TestVersionArg(t *testing.T) {
	t.Parallel()
	for _, s := range []string{"version", "-version", "--version"} {
		if !versionArg(s) {
			t.Fatalf("versionArg(%q) want true", s)
		}
	}
	for _, s := range []string{"", "v", "serve", "-v"} {
		if versionArg(s) {
			t.Fatalf("versionArg(%q) want false", s)
		}
	}
}

func TestRun_unknownSubcommand(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "nope"}); g != 1 {
		t.Fatalf("exit %d want 1", g)
	}
}

func TestRun_serveExtraArgs(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "serve", "extra"}); g != 1 {
		t.Fatalf("exit %d want 1", g)
	}
}

func TestRun_userWithoutSubcommand(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "user"}); g != 2 {
		t.Fatalf("exit %d want 2", g)
	}
}

func TestRun_userUnknownSubcommand(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "user", "nope"}); g != 2 {
		t.Fatalf("exit %d want 2", g)
	}
}

func TestRun_userAdd_missingEmailPassword(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "user", "add"}); g != 2 {
		t.Fatalf("exit %d want 2", g)
	}
}

func TestRun_userPassword_missingEmailPassword(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "user", "password"}); g != 2 {
		t.Fatalf("exit %d want 2", g)
	}
}

func TestRun_userAdd_invalidRole(t *testing.T) {
	t.Parallel()
	if g := run([]string{"moana", "user", "add", "-email", "a@b.co", "-password", "secret", "-role", "superadmin"}); g != 2 {
		t.Fatalf("exit %d want 2", g)
	}
}

func TestRun_userPassword_unknownEmail(t *testing.T) {
	// Do not use t.Parallel: MOANA_DB_PATH must not race with other tests.
	t.Setenv("MOANA_DB_PATH", ":memory:")
	if g := run([]string{"moana", "user", "password", "-email", "nobody@example.com", "-password", "secret"}); g != 1 {
		t.Fatalf("exit %d want 1 (user not in empty DB)", g)
	}
}
