package store

import "strings"

// sqliteUniqueError reports SQLite UNIQUE constraint failures from modernc.org/sqlite.
func sqliteUniqueError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
