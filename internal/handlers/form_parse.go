package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

// categoryIDFromForm parses a positive category id from a form field string, or 0 if missing/invalid.
// Leading/trailing space is trimmed (same rules as formPositiveInt64).
func categoryIDFromForm(v string) int64 {
	id, _ := positiveInt64String(v)
	return id
}

// formPositiveInt64 parses FormValue(name) as int64 > 0, or false if missing/invalid.
func formPositiveInt64(r *http.Request, name string) (int64, bool) {
	return positiveInt64String(r.FormValue(name))
}

// positiveInt64String parses v as a positive base-10 int64 after TrimSpace.
func positiveInt64String(v string) (int64, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
