package handlers

import (
	"errors"
	"net/http"
	"time"

	"moana/internal/httperr"
	"moana/internal/store"
	"moana/internal/tz"
)

// Notifications renders GET /notifications (user-scoped inbox from the store).
func (a *App) Notifications(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	items, err := a.Store.ListNotificationsForUser(ctx, u.ID, 0)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	loc := tz.DisplayLocation(r)
	p := NotificationsPageData{
		MarkReadPath: NotificationsMarkReadPath,
		Items:        make([]NotificationView, 0, len(items)),
	}
	for _, n := range items {
		when := n.CreatedAt.In(loc).Format(time.DateTime)
		p.Items = append(p.Items, NotificationView{
			ID:     n.ID,
			Body:   n.Body,
			When:   when,
			Unread: n.ReadAt == nil,
		})
	}
	a.renderShell(w, "notifications.html", p, layoutShell("Notifications", "notifications", u))
}

// NotificationMarkRead handles POST /notifications/read (mark one row read for the signed-in user).
func (a *App) NotificationMarkRead(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseForm(w, r) {
		return
	}
	id, ok := formPositiveInt64(r, NotificationFieldID)
	if !ok {
		http.Redirect(w, r, NotificationsPath, http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	if err := a.Store.MarkNotificationRead(ctx, u.ID, id); err != nil {
		if errors.Is(err, store.ErrNotificationNotFound) {
			http.NotFound(w, r)
			return
		}
		httperr.Internal(w, r, err)
		return
	}
	http.Redirect(w, r, NotificationsPath, http.StatusSeeOther)
}

// NotificationsPageData is template data for notifications.html.
type NotificationsPageData struct {
	MarkReadPath string
	Items        []NotificationView
}

// NotificationView is one row in the notifications template.
type NotificationView struct {
	ID     int64
	Body   string
	When   string
	Unread bool
}
