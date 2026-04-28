package store

import (
	"database/sql"
	"fmt"
)

// execExactlyOneRow verifies a successful Exec affected exactly one row.
// errZero is returned when RowsAffected is 0 (e.g. [sql.ErrNoRows] or [ErrCategoryNotFound]).
// Use only after ExecContext returned a nil error.
func execExactlyOneRow(res sql.Result, errZero error) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errZero
	}
	if n != 1 {
		return fmt.Errorf("expected one row affected, got %d", n)
	}
	return nil
}
