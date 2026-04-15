package handlers

import (
	"net/http"
	"strconv"
)

// pathPositiveInt64 parses r.PathValue(name) as a positive int64 id, or returns 0, false.
func pathPositiveInt64(r *http.Request, name string) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}
