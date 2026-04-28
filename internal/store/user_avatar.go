package store

import (
	"context"
	"database/sql"
)

// IncrementUserAvatarRev bumps the stored avatar version after a new image file is written.
func (s *Store) IncrementUserAvatarRev(ctx context.Context, userID int64) error {
	res, err := s.DB.ExecContext(ctx, sqlUserIncrementAvatarRev, userID)
	if err != nil {
		return err
	}
	return execExactlyOneRow(res, sql.ErrNoRows)
}

// ClearUserAvatar removes the avatar revision (and the handler should delete the on-disk file).
func (s *Store) ClearUserAvatar(ctx context.Context, userID int64) error {
	res, err := s.DB.ExecContext(ctx, sqlUserClearAvatar, userID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}
