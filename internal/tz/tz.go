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
	return requestLocation(r).String()
}

// DisplayLocation is the browser time zone for this request, or UTC.
// The result is never nil.
func DisplayLocation(r *http.Request) *time.Location {
	return requestLocation(r)
}

// requestLocation resolves moana_tz once per call ([timeutil.LoadLocation] maps invalid names to UTC).
func requestLocation(r *http.Request) *time.Location {
	if r == nil {
		return time.UTC
	}
	c, err := r.Cookie(CookieName)
	if err != nil || c.Value == "" {
		return time.UTC
	}
	v := strings.TrimSpace(c.Value)
	if v == "" {
		return time.UTC
	}
	return timeutil.LoadLocation(v)
}
