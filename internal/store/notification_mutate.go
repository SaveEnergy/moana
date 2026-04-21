package store

import (
	"context"
	"errors"

	"moana/internal/timeutil"
)

// ErrNotificationNotFound is returned when no notification row matches id + user ([MarkNotificationRead]).
var ErrNotificationNotFound = errors.New("notification not found")

// MarkNotificationRead sets read_at for a notification owned by userID.
func (s *Store) MarkNotificationRead(ctx context.Context, userID, notificationID int64) error {
	now := timeutil.NowSQLiteUTC()
	res, err := s.DB.ExecContext(ctx, sqlNotificationMarkRead, now, notificationID, userID)
	if err != nil {
		return err
	}
	return execExactlyOneRow(res, ErrNotificationNotFound)
}
