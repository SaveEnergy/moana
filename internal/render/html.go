package render

import "net/http"

// writeHTML sets the standard HTML Content-Type and writes the body.
func writeHTML(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(body)
}
