package handlers

import (
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
	p := NotificationsPageData{Items: make([]NotificationView, 0, len(items))}
	for _, n := range items {
		when := n.CreatedAt.In(loc).Format(time.DateTime)
		p.Items = append(p.Items, NotificationView{
			Body:   n.Body,
			When:   when,
			Unread: n.ReadAt == nil,
		})
	}
	a.renderShell(w, "notifications.html", p, layoutShell("Notifications", "notifications", u))
}

// NotificationsPageData is template data for notifications.html.
type NotificationsPageData struct {
	Items []NotificationView
}

// NotificationView is one row in the notifications template.
type NotificationView struct {
	Body   string
	When   string
	Unread bool
}
