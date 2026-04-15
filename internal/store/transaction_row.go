package store

import (
	"database/sql"

	"moana/internal/timeutil"
)

func hydrateTransaction(t *Transaction, occ, cre string, catID sql.NullInt64) error {
	t.CategoryID = catID
	var err error
	t.OccurredAt, err = timeutil.ParseSQLiteTimestamp(occ)
	if err != nil {
		return err
	}
	t.CreatedAt, err = timeutil.ParseSQLiteTimestamp(cre)
	return err
}
