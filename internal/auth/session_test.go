package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignReadSessionRoundTrip(t *testing.T) {
	t.Parallel()
	secret := []byte("test-hmac-secret-at-least-32-bytes-long-ok")
	w := httptest.NewRecorder()
	err := SignSession(w, secret, SessionPayload{UserID: 42, Role: "admin"}, time.Hour, false)
	if err != nil {
		t.Fatal(err)
	}
	var raw string
	for _, c := range w.Result().Cookies() {
		if c.Name == "moana_session" {
			raw = c.Value
			break
		}
	}
	if raw == "" {
		t.Fatal("no session cookie")
	}
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "moana_session", Value: raw})
	p, err := ReadSession(req, secret)
	if err != nil {
		t.Fatal(err)
	}
	if p.UserID != 42 || p.Role != "admin" {
		t.Fatalf("payload %+v", p)
	}
}
