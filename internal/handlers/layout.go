package handlers

import (
	"net/http"
	"time"

	"moana/internal/render"
	"moana/internal/store"
)

// LayoutData is the authenticated page shell; see [render.LayoutData].
type LayoutData = render.LayoutData

// layoutShell builds standard authenticated shell metadata (title, nav highlight, footer year).
func layoutShell(title, navKey string, u *store.User) LayoutData {
	return LayoutData{
		Title:  title,
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: navKey,
	}
}

// layoutShellMain is like layoutShell but sets MainClass (e.g. settings-shell on the main column).
func layoutShellMain(title, navKey, mainClass string, u *store.User) LayoutData {
	return LayoutData{
		Title:     title,
		User:      u,
		Year:      time.Now().UTC().Year(),
		Active:    navKey,
		MainClass: mainClass,
	}
}

// renderShell executes the named page template (e.g. dashboard.html) into the layout body.
func (a *App) renderShell(w http.ResponseWriter, contentTemplate string, data any, ld LayoutData) {
	a.Render.Shell(w, contentTemplate, data, ld, a.Config.RepoURL)
}

// renderSimple executes a standalone template (e.g. login.html) without the app shell.
func (a *App) renderSimple(w http.ResponseWriter, name string, data any) {
	a.Render.Simple(w, name, data)
}
