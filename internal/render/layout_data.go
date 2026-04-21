package render

import (
	"html/template"

	"moana/internal/store"
)

// LayoutData is the outer shell for authenticated pages (layout.html).
type LayoutData struct {
	Title  string
	User   *store.User
	Body   template.HTML
	Year   int
	Active string
	// MainClass is an optional extra class on <main> (e.g. settings-shell).
	MainClass string
	// UnreadNotificationCount is inbox unread rows for the shell badge (see [handlers.App] renderShell).
	UnreadNotificationCount int64
	// RepoURL is set by [Engine.Shell] from config (public GitHub / source link).
	RepoURL string
}
