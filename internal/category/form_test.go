package category

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseColorFromForm_preset(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("POST", "/", strings.NewReader("color=%23ea580c"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if got := ParseColorFromForm(r); got != "#ea580c" {
		t.Fatalf("got %q", got)
	}
}

func TestParseColorFromForm_custom(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("POST", "/", strings.NewReader("color=custom&color_custom=%23AaBbCc"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if got := ParseColorFromForm(r); got != "#aabbcc" {
		t.Fatalf("got %q", got)
	}
}

func TestParseColorFromForm_invalid(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("POST", "/", strings.NewReader("color=not-a-hex"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if got := ParseColorFromForm(r); got != "" {
		t.Fatalf("got %q", got)
	}
}
