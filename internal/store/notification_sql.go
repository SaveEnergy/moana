package store

const sqlNotificationListForUser = `
SELECT id, body, read_at, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`

const sqlNotificationInsert = `INSERT INTO notifications (user_id, body, read_at, created_at) VALUES (?, ?, NULL, ?)`
