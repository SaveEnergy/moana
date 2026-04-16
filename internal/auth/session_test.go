package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func sessionCookieValue(w *httptest.ResponseRecorder) string {
	for _, c := range w.Result().Cookies() {
		if c.Name == cookieName {
			return c.Value
		}
	}
	return ""
}

func TestSignReadSessionRoundTrip(t *testing.T) {
	t.Parallel()
	secret := []byte("test-hmac-secret-at-least-32-bytes-long-ok")
	w := httptest.NewRecorder()
	err := SignSession(w, secret, SessionPayload{UserID: 42, Role: "admin"}, time.Hour, false)
	if err != nil {
		t.Fatal(err)
	}
	raw := sessionCookieValue(w)
	if raw == "" {
		t.Fatal("no session cookie")
	}
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: raw})
	p, err := ReadSession(req, secret)
	if err != nil {
		t.Fatal(err)
	}
	if p.UserID != 42 || p.Role != "admin" {
		t.Fatalf("payload %+v", p)
	}
}

func TestReadSession_expired(t *testing.T) {
	t.Parallel()
	secret := []byte("test-hmac-secret-at-least-32-bytes-long-ok")
	w := httptest.NewRecorder()
	if err := SignSession(w, secret, SessionPayload{UserID: 1, Role: "user"}, -time.Hour, false); err != nil {
		t.Fatal(err)
	}
	raw := sessionCookieValue(w)
	if raw == "" {
		t.Fatal("no session cookie")
	}
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: raw})
	if _, err := ReadSession(req, secret); err == nil {
		t.Fatal("expected expired session error")
	}
}

func TestReadSession_wrongSecret(t *testing.T) {
	t.Parallel()
	secret := []byte("test-hmac-secret-at-least-32-bytes-long-ok")
	w := httptest.NewRecorder()
	if err := SignSession(w, secret, SessionPayload{UserID: 1, Role: "user"}, time.Hour, false); err != nil {
		t.Fatal(err)
	}
	raw := sessionCookieValue(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: raw})
	if _, err := ReadSession(req, []byte("different-secret----not-same-as-above-ok")); err == nil {
		t.Fatal("expected error for wrong HMAC secret")
	}
}

func TestClearSession_expiresCookie(t *testing.T) {
	t.Parallel()
	t.Run("secure", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		ClearSession(w, true)
		cs := w.Result().Cookies()
		if len(cs) != 1 {
			t.Fatalf("cookies: %d", len(cs))
		}
		c := cs[0]
		if c.Name != cookieName || c.Value != "" || c.MaxAge != -1 || !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode {
			t.Fatalf("cookie: %+v", c)
		}
	})
	t.Run("not_secure", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		ClearSession(w, false)
		c := w.Result().Cookies()[0]
		if c.Secure {
			t.Fatal("want Secure false")
		}
	})
}

func TestReadSession_invalidRole(t *testing.T) {
	t.Parallel()
	secret := []byte("test-hmac-secret-at-least-32-bytes-long-ok")
	w := httptest.NewRecorder()
	if err := SignSession(w, secret, SessionPayload{UserID: 1, Role: "guest"}, time.Hour, false); err != nil {
		t.Fatal(err)
	}
	raw := sessionCookieValue(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: raw})
	if _, err := ReadSession(req, secret); err == nil {
		t.Fatal("expected error for invalid role")
	}
}
