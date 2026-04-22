package handlers

import (
	"moana/internal/config"
	"moana/internal/mail"
	"moana/internal/render"
	"moana/internal/store"
)

// App wires HTTP handlers to shared dependencies. Prefer moana/internal/app.New for a fully wired instance.
type App struct {
	Config *config.Config
	Store  *store.Store
	Render *render.Engine
	// PasswordReset sends “forgot password” email when SendGrid is configured; nil otherwise.
	PasswordReset mail.PasswordResetSender
	// AvatarDir is the on-disk directory for user profile JPEGs (see [App.UserAvatarGet]).
	AvatarDir string
}
