package store

import "strings"

// normalizeUserEmail trims surrounding ASCII whitespace for lookups and persistence.
func normalizeUserEmail(email string) string {
	return strings.TrimSpace(email)
}
