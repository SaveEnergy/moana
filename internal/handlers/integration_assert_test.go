package handlers_test

import (
	"strings"
	"testing"
)

// assertBodyHasErrorAlert asserts the HTML includes an error alert.
// Some templates use `class="alert alert-error"`; settings uses `alert alert-error settings-alert` —
// matching the common prefix avoids brittle full-string checks.
func assertBodyHasErrorAlert(t *testing.T, body string) {
	t.Helper()
	if !strings.Contains(body, `class="alert alert-error`) {
		t.Fatal("expected error alert (class prefix alert-error)")
	}
}
