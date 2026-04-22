package mail

import (
	"context"
	"strings"
)

// PasswordResetSender is implemented by [sendGridSender] for dependency injection in [handlers.App].
type PasswordResetSender interface {
	SendPasswordReset(ctx context.Context, toEmail, resetURL string) error
}

// PublicResetURL builds an absolute https? URL for a path starting with /.
func PublicResetURL(publicBase, path string) string {
	return strings.TrimRight(publicBase, "/") + path
}
