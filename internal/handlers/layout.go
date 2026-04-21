package handlers

import (
	"net/http"
	"time"

	"moana/internal/httperr"
	"moana/internal/render"
	"moana/internal/store"
)

// LayoutData is the authenticated page shell; see [render.LayoutData].
type LayoutData = render.LayoutData

// shellYear is the footer copyright year for authenticated shells (single clock read per call).
func shellYear() int {
	return time.Now().UTC().Year()
}

// layoutShell builds standard authenticated shell metadata (title, nav highlight, footer year).
func layoutShell(title, navKey string, u *store.User) LayoutData {
	return layoutData(title, navKey, "", u)
}

// layoutShellMain is like layoutShell but sets MainClass (e.g. settings-shell on the main column).
func layoutShellMain(title, navKey, mainClass string, u *store.User) LayoutData {
	return layoutData(title, navKey, mainClass, u)
}

func layoutData(title, navKey, mainClass string, u *store.User) LayoutData {
	return layoutDataWithUnread(title, navKey, mainClass, u, 0)
}

// layoutDataWithUnread builds shell metadata; unread is used for the top bar notification badge.
func layoutDataWithUnread(title, navKey, mainClass string, u *store.User, unread int64) LayoutData {
	return LayoutData{
		Title:                   title,
		User:                    u,
		Year:                    shellYear(),
		Active:                  navKey,
		MainClass:               mainClass,
		UnreadNotificationCount: unread,
	}
}

// renderShell executes the named page template (e.g. dashboard.html) into the layout body.
// Unread count comes from [WithAuth] request context when present; otherwise it falls back to
// [store.Store.CountUnreadNotificationsForUser].
func (a *App) renderShell(w http.ResponseWriter, r *http.Request, contentTemplate string, pageData any, title, navKey, mainClass string, u *store.User) {
	var unread int64
	if n, ok := unreadNotificationCountFromContext(r); ok {
		unread = n
	} else {
		var err error
		unread, err = a.Store.CountUnreadNotificationsForUser(r.Context(), u.ID)
		if err != nil {
			httperr.Internal(w, r, err)
			return
		}
	}
	a.Render.Shell(w, contentTemplate, pageData, layoutDataWithUnread(title, navKey, mainClass, u, unread), a.Config.RepoURL)
}

// renderSimple executes a standalone template (e.g. login.html) without the app shell.
func (a *App) renderSimple(w http.ResponseWriter, name string, data any) {
	a.Render.Simple(w, name, data)
}
