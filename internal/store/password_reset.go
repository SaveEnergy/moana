package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"moana/internal/timeutil"
)

// ErrPasswordResetInvalid means the token is unknown, expired, or already used.
var ErrPasswordResetInvalid = errors.New("password reset link is invalid or expired")

// ReplacePasswordResetToken replaces any existing token for the user, inserts a new one.
func (s *Store) ReplacePasswordResetToken(ctx context.Context, userID int64, tokenHash []byte, expires, created time.Time) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, sqlPasswordResetDeleteByUser, userID); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sqlPasswordResetInsert, userID, tokenHash,
		timeutil.FormatSQLiteUTC(expires), timeutil.FormatSQLiteUTC(created))
	if err != nil {
		return err
	}
	return tx.Commit()
}

// RedeemPasswordResetToken sets the user's password and removes all password-reset tokens for that user.
func (s *Store) RedeemPasswordResetToken(ctx context.Context, tokenHash []byte, newPasswordHash []byte) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var (
		uid  int64
		expS string
	)
	err = tx.QueryRowContext(ctx, sqlPasswordResetSelect, tokenHash).Scan(&uid, &expS)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPasswordResetInvalid
	}
	if err != nil {
		return err
	}
	exp, err := timeutil.ParseSQLiteTimestampUTC(expS)
	if err != nil {
		return err
	}
	if time.Now().UTC().After(exp) {
		_, _ = tx.ExecContext(ctx, sqlPasswordResetDeleteByHash, tokenHash)
		return ErrPasswordResetInvalid
	}
	res, err := tx.ExecContext(ctx, sqlUserUpdatePassword, newPasswordHash, uid)
	if err != nil {
		return err
	}
	if err := execExactlyOneRow(res, sql.ErrNoRows); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, sqlPasswordResetDeleteByUserID, uid); err != nil {
		return err
	}
	return tx.Commit()
}
