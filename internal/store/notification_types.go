package store

import "time"

// Notification is one in-app notification row for a user.
type Notification struct {
	ID        int64
	Body      string
	ReadAt    *time.Time
	CreatedAt time.Time
}
