package htmlview

import "testing"

func TestUserInitial(t *testing.T) {
	t.Parallel()
	if got := UserInitial(""); got != "?" {
		t.Fatalf("empty: got %q", got)
	}
	if got := UserInitial("  alice@x.test  "); got != "A" {
		t.Fatalf("got %q", got)
	}
}

func TestAdminDisplayName(t *testing.T) {
	t.Parallel()
	if got := AdminDisplayName("john.doe@example.com"); got != "John Doe" {
		t.Fatalf("got %q", got)
	}
	if got := AdminDisplayName("  "); got != "" {
		t.Fatalf("whitespace: got %q", got)
	}
}

func TestProfileDisplayName(t *testing.T) {
	t.Parallel()
	if got := ProfileDisplayName("Jane", "Doe", "ignore@test"); got != "Jane Doe" {
		t.Fatalf("got %q", got)
	}
	if got := ProfileDisplayName("", "", "a.b@x.com"); got != "A B" {
		t.Fatalf("fallback: got %q", got)
	}
}

func TestProfileInitial(t *testing.T) {
	t.Parallel()
	if got := ProfileInitial("Jane", "", ""); got != "J" {
		t.Fatalf("first: got %q", got)
	}
	if got := ProfileInitial("", "Doe", ""); got != "D" {
		t.Fatalf("last: got %q", got)
	}
	if got := ProfileInitial("", "", "zoe@x.test"); got != "Z" {
		t.Fatalf("email: got %q", got)
	}
}
