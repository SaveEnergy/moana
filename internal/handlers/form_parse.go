package handlers

import "strconv"

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
