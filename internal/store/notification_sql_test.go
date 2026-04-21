package store

import "testing"

func TestSqlNotificationListForUser_stable(t *testing.T) {
	t.Parallel()
	want := `
SELECT id, body, read_at, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`
	if sqlNotificationListForUser != want {
		t.Fatalf("sqlNotificationListForUser drift")
	}
}

func TestSqlNotificationInsert_stable(t *testing.T) {
	t.Parallel()
	want := `INSERT INTO notifications (user_id, body, read_at, created_at) VALUES (?, ?, NULL, ?)`
	if sqlNotificationInsert != want {
		t.Fatalf("sqlNotificationInsert drift")
	}
}

func TestSqlNotificationMarkRead_stable(t *testing.T) {
	t.Parallel()
	want := `UPDATE notifications SET read_at = ? WHERE id = ? AND user_id = ?`
	if sqlNotificationMarkRead != want {
		t.Fatalf("sqlNotificationMarkRead drift")
	}
}

func TestSqlNotificationCountUnreadForUser_stable(t *testing.T) {
	t.Parallel()
	want := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read_at IS NULL`
	if sqlNotificationCountUnreadForUser != want {
		t.Fatalf("sqlNotificationCountUnreadForUser drift")
	}
}
