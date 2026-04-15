package tz

import (
	"net/http"
	"testing"
	"time"
)

func TestCookieZone_nilRequest(t *testing.T) {
	t.Parallel()
	if got := CookieZone(nil); got != "UTC" {
		t.Fatalf("got %q", got)
	}
}

func TestCookieZone_missing(t *testing.T) {
	t.Parallel()
	r := &http.Request{Header: http.Header{}}
	if got := CookieZone(r); got != "UTC" {
		t.Fatalf("got %q", got)
	}
}

func TestCookieZone_valid(t *testing.T) {
	t.Parallel()
	r := &http.Request{Header: http.Header{}}
	r.AddCookie(&http.Cookie{Name: CookieName, Value: "Europe/Berlin"})
	if got := CookieZone(r); got != "Europe/Berlin" {
		t.Fatalf("got %q", got)
	}
}

func TestCookieZone_invalidIANA(t *testing.T) {
	t.Parallel()
	r := &http.Request{Header: http.Header{}}
	r.AddCookie(&http.Cookie{Name: CookieName, Value: "Not/A/Real/Zone"})
	if got := CookieZone(r); got != "UTC" {
		t.Fatalf("got %q", got)
	}
}

func TestDisplayLocation(t *testing.T) {
	t.Parallel()
	r := &http.Request{Header: http.Header{}}
	r.AddCookie(&http.Cookie{Name: CookieName, Value: "America/New_York"})
	loc := DisplayLocation(r)
	if loc.String() != "America/New_York" {
		t.Fatalf("got %v", loc)
	}
	if _, err := time.LoadLocation("America/New_York"); err != nil {
		t.Fatal(err)
	}
}
