package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

// categoryIDFromForm parses a positive category id from a form field, or 0 if missing/invalid.
func categoryIDFromForm(v string) int64 {
	if v == "" {
		return 0
	}
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil || id <= 0 {
		return 0
	}
	return id
}

// formPositiveInt64 parses TrimSpace(FormValue(name)) as int64 > 0, or false if missing/invalid.
func formPositiveInt64(r *http.Request, name string) (int64, bool) {
	v := strings.TrimSpace(r.FormValue(name))
	if v == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
