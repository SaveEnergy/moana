package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRedirectSettingsErr_roundtripQuery(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, SettingsPath, nil)
	redirectSettingsErr(rec, req, `bad & "quotes"`)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("code %d", rec.Code)
	}
	loc := rec.Result().Header.Get("Location")
	u, err := url.Parse(loc)
	if err != nil {
		t.Fatal(err)
	}
	if u.Path != SettingsPath {
		t.Fatalf("path %q", u.Path)
	}
	if got := u.Query().Get(SettingsErrorQueryParam); got != `bad & "quotes"` {
		t.Fatalf("err param %q", got)
	}
}

func TestRedirectSettingsOK_roundtripQuery(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, SettingsPath, nil)
	redirectSettingsOK(rec, req, "ok&key")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("code %d", rec.Code)
	}
	loc := rec.Result().Header.Get("Location")
	u, err := url.Parse(loc)
	if err != nil {
		t.Fatal(err)
	}
	if got := u.Query().Get(SettingsOKQueryParam); got != "ok&key" {
		t.Fatalf("ok param %q", got)
	}
}
