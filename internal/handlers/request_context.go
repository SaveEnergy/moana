package handlers

import (
	"context"
	"net/http"
)

type ctxKeyUnreadNotif int

const keyUnreadNotif ctxKeyUnreadNotif = 0

func withUnreadNotificationCount(r *http.Request, n int64) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), keyUnreadNotif, n))
}

func unreadNotificationCountFromContext(r *http.Request) (int64, bool) {
	v := r.Context().Value(keyUnreadNotif)
	if v == nil {
		return 0, false
	}
	n, ok := v.(int64)
	return n, ok
}
