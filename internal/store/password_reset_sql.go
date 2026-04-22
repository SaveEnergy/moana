package store

// Password reset: short-lived, single-use (row deleted on redeem).
const (
	sqlPasswordResetDeleteByUser   = `DELETE FROM password_reset_tokens WHERE user_id = ?`
	sqlPasswordResetInsert         = `INSERT INTO password_reset_tokens (user_id, token_hash, expires_at, created_at) VALUES (?, ?, ?, ?)`
	sqlPasswordResetSelect         = `SELECT user_id, expires_at FROM password_reset_tokens WHERE token_hash = ?`
	sqlPasswordResetDeleteByUserID = `DELETE FROM password_reset_tokens WHERE user_id = ?`
	sqlPasswordResetDeleteByHash   = `DELETE FROM password_reset_tokens WHERE token_hash = ?`
)
