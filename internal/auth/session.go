package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// SessionPayload is the signed cookie content.
type SessionPayload struct {
	UserID int64  `json:"uid"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

const cookieName = "moana_session"

// SignSession sets a signed session cookie.
func SignSession(w http.ResponseWriter, secret []byte, p SessionPayload, maxAge time.Duration, secure bool) error {
	p.Exp = time.Now().Add(maxAge).Unix()
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	sig := mac.Sum(nil)
	token := base64.RawURLEncoding.EncodeToString(body) + "." + base64.RawURLEncoding.EncodeToString(sig)

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

// ClearSession removes the session cookie.
func ClearSession(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ReadSession verifies and parses the session cookie.
func ReadSession(r *http.Request, secret []byte) (*SessionPayload, error) {
	c, err := r.Cookie(cookieName)
	if err != nil || c.Value == "" {
		return nil, errors.New("no session")
	}
	parts := strings.Split(c.Value, ".")
	if len(parts) != 2 {
		return nil, errors.New("bad session format")
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return nil, errors.New("bad signature")
	}
	var p SessionPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	if time.Now().Unix() > p.Exp {
		return nil, errors.New("expired")
	}
	if p.UserID <= 0 || (p.Role != "user" && p.Role != "admin") {
		return nil, errors.New("invalid session")
	}
	return &p, nil
}
