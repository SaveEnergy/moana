package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithUnreadNotificationCount_contextRoundTrip(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r2 := withUnreadNotificationCount(r, 42)
	n, ok := unreadNotificationCountFromContext(r2)
	if !ok || n != 42 {
		t.Fatalf("got ok=%v n=%d", ok, n)
	}
	if _, ok := unreadNotificationCountFromContext(r); ok {
		t.Fatal("expected original request without unread key")
	}
}
