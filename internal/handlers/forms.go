package handlers

import "net/http"

// requireParseForm parses the request body for URL-encoded forms. On failure it sends 400 and returns false.
func requireParseForm(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return false
	}
	return true
}

// requireParseFormSettings is like requireParseForm but redirects to /settings with an error message (settings flows).
func requireParseFormSettings(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		redirectSettingsErr(w, r, "Invalid form.")
		return false
	}
	return true
}
