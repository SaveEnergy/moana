package handlers

import (
	"net/http/httptest"
	"testing"
)

func TestPathPositiveInt64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   string
		wantID  int64
		wantOK  bool
	}{
		{"valid", "42", 42, true},
		{"valid_max_reasonable", "9223372036854775807", 1<<63 - 1, true},
		{"empty", "", 0, false},
		{"zero", "0", 0, false},
		{"negative", "-1", 0, false},
		{"not_a_number", "abc", 0, false},
		{"float", "3.14", 0, false},
		// strconv.ParseInt accepts an optional leading + (same as a positive id).
		{"leading_plus", "+5", 5, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest("GET", "/x", nil)
			r.SetPathValue("id", tt.value)
			got, ok := pathPositiveInt64(r, "id")
			if ok != tt.wantOK || got != tt.wantID {
				t.Fatalf("pathPositiveInt64: got (%d, %v), want (%d, %v)", got, ok, tt.wantID, tt.wantOK)
			}
		})
	}
}

func TestPathPositiveInt64_wrongName(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "/x", nil)
	r.SetPathValue("id", "7")
	if id, ok := pathPositiveInt64(r, "other"); ok || id != 0 {
		t.Fatalf("got (%d, %v), want (0, false)", id, ok)
	}
}
