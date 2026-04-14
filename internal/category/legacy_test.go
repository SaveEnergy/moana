package category

import "testing"

func TestNormalizeStoredIcon(t *testing.T) {
	t.Parallel()
	if got := NormalizeStoredIcon("  "); got != "" {
		t.Fatalf("empty: %q", got)
	}
	if got := NormalizeStoredIcon("wallet"); got != "wallet" {
		t.Fatalf("valid id: %q", got)
	}
	if got := NormalizeStoredIcon("💰"); got != "wallet" {
		t.Fatalf("legacy emoji: %q", got)
	}
}
