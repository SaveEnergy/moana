package store

import (
	"errors"
	"testing"
)

func TestSqliteUniqueError(t *testing.T) {
	t.Parallel()
	if !sqliteUniqueError(errors.New(`constraint failed: UNIQUE constraint failed: categories.household_id, categories.name`)) {
		t.Fatal("expected true for UNIQUE substring")
	}
	if sqliteUniqueError(errors.New("no such table: foo")) {
		t.Fatal("expected false for unrelated error")
	}
	if sqliteUniqueError(nil) {
		t.Fatal("expected false for nil")
	}
}
