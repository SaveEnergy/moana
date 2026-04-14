package safepath

import "strings"

// Default is used when the requested next URL is missing or unsafe.
const Default = "/history"

// Internal returns next if it is a same-origin path, otherwise Default.
func Internal(next string) string {
	next = strings.TrimSpace(next)
	if next == "" || !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return Default
	}
	return next
}
