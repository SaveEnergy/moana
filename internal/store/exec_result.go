package store

import "database/sql"

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
	return nil
}
