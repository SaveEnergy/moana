package store

import (
	"database/sql"
	"errors"

	"moana/internal/timeutil"
)

func scanUser(row *sql.Row) (*User, error) {
	var u User
	var created string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &created,
		&u.HouseholdID, &u.FirstName, &u.LastName, &u.HouseholdRole)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t, err := timeutil.ParseSQLiteTimestamp(created)
	if err != nil {
		return nil, err
	}
	u.CreatedAt = t.UTC()
	return &u, nil
}
