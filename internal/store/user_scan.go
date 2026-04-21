package store

import (
	"database/sql"
	"errors"

	"moana/internal/timeutil"
)

func finalizeUserCreatedAt(u *User, created string) error {
	t, err := timeutil.ParseSQLiteTimestampUTC(created)
	if err != nil {
		return err
	}
	u.CreatedAt = t
	return nil
}

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
	if err := finalizeUserCreatedAt(&u, created); err != nil {
		return nil, err
	}
	return &u, nil
}

// scanUserAndUnread scans a user row plus an unread notification count (see [sqlUserGetByIDWithUnreadCount]).
func scanUserAndUnread(row *sql.Row) (*User, int64, error) {
	var u User
	var created string
	var unread int64
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &created,
		&u.HouseholdID, &u.FirstName, &u.LastName, &u.HouseholdRole, &unread)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	if err := finalizeUserCreatedAt(&u, created); err != nil {
		return nil, 0, err
	}
	return &u, unread, nil
}
