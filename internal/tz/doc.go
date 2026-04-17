// Package tz resolves the browser IANA time zone from the moana_tz cookie (set by the frontend).
// [DisplayLocation] always returns a non-nil *Location (UTC when the cookie is missing/invalid).
// [CookieZone] and [DisplayLocation] are covered in tz_test.go.
package tz
