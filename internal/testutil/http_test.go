package testutil

import "testing"

func TestNewCookieClient_hasJar(t *testing.T) {
	t.Parallel()
	c := NewCookieClient(t)
	if c.Jar == nil {
		t.Fatal("expected non-nil cookie jar for session integration tests")
	}
}
