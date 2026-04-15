package store

import "strings"

// escapeSQLLikePattern escapes %, _, and the escape character for SQLite LIKE ... ESCAPE '!'.
// See https://www.sqlite.org/lang_expr.html — without this, user search terms are not literal.
func escapeSQLLikePattern(s string) string {
	s = strings.ReplaceAll(s, "!", "!!")
	s = strings.ReplaceAll(s, "%", "!%")
	s = strings.ReplaceAll(s, "_", "!_")
	return s
}
