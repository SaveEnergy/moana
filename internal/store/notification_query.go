package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"moana/internal/timeutil"
)

// ErrInvalidNotificationBody is returned when [Store.InsertNotification] gets an empty body after trim.
var ErrInvalidNotificationBody = errors.New("notification body must not be empty")

// defaultNotificationListLimit caps list size when callers pass non-positive limits.
const defaultNotificationListLimit = 50

// maxNotificationListLimit avoids unbounded reads.
const maxNotificationListLimit = 500

// ListNotificationsForUser returns notifications for the user, newest first.
func (s *Store) ListNotificationsForUser(ctx context.Context, userID int64, limit int) ([]Notification, error) {
	if limit < 1 {
		limit = defaultNotificationListLimit
	}
	if limit > maxNotificationListLimit {
		limit = maxNotificationListLimit
	}
	rows, err := s.DB.QueryContext(ctx, sqlNotificationListForUser, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Notification, 0, min(limit, 64))
	for rows.Next() {
		var n Notification
		var readAt sql.NullString
		var created string
		if err := rows.Scan(&n.ID, &n.Body, &readAt, &created); err != nil {
			return nil, err
		}
		if readAt.Valid {
			t, err := timeutil.ParseSQLiteTimestamp(readAt.String)
			if err != nil {
				return nil, err
			}
			tt := t.UTC()
			n.ReadAt = &tt
		}
		var err error
		n.CreatedAt, err = timeutil.ParseSQLiteTimestamp(created)
		if err != nil {
			return nil, err
		}
		n.CreatedAt = n.CreatedAt.UTC()
		out = append(out, n)
	}
	return out, rows.Err()
}

// InsertNotification stores a notification for userID. Body is trimmed; empty body returns [ErrInvalidNotificationBody].
func (s *Store) InsertNotification(ctx context.Context, userID int64, body string) (int64, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return 0, ErrInvalidNotificationBody
	}
	now := timeutil.NowSQLiteUTC()
	res, err := s.DB.ExecContext(ctx, sqlNotificationInsert, userID, body, now)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
