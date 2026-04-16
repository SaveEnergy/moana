package tz

import (
	"net/http"
	"strings"
	"time"

	"moana/internal/timeutil"
)

// CookieName is the cookie set by the client with the browser's IANA zone (see internal/assets/static/js/app.js).
const CookieName = "moana_tz"

// CookieZone returns an IANA zone name from the cookie, or "UTC" if missing/invalid.
func CookieZone(r *http.Request) string {
	if r == nil {
		return "UTC"
	}
	c, err := r.Cookie(CookieName)
	if err != nil || c.Value == "" {
		return "UTC"
	}
	v := strings.TrimSpace(c.Value)
	if _, err := time.LoadLocation(v); err != nil {
		return "UTC"
	}
	return v
}

// DisplayLocation is the browser time zone for this request, or UTC.
func DisplayLocation(r *http.Request) *time.Location {
	return timeutil.LoadLocation(CookieZone(r))
}
