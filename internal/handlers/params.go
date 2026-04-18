package handlers

import "net/http"

// pathPositiveInt64 parses r.PathValue(name) as a positive int64 id, or returns 0, false.
// Leading/trailing space is trimmed (same rules as formPositiveInt64).
func pathPositiveInt64(r *http.Request, name string) (int64, bool) {
	return positiveInt64String(r.PathValue(name))
}
