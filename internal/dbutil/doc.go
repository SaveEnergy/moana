// Package dbutil provides helpers for opening SQLite and constructing [moana/internal/store.Store].
// [MustOpenMemStore] is for unit tests; production and CLI use [OpenStore] or [moana/internal/db.Open].
// open_test.go smoke-tests [OpenStore] on :memory:.
package dbutil
