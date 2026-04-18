package historyview

import "testing"

func TestPartialDateFilter(t *testing.T) {
	t.Parallel()
	if !partialDateFilter("2026-01-01", "") {
		t.Fatal("expected true when only from is set")
	}
	if !partialDateFilter("", "2026-01-31") {
		t.Fatal("expected true when only to is set")
	}
	if partialDateFilter("", "") {
		t.Fatal("expected false when both empty")
	}
	if partialDateFilter("2026-01-01", "2026-01-31") {
		t.Fatal("expected false when both set")
	}
}
