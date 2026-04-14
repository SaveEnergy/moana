package handlers

import (
	"net/http"
	"time"

	"moana/internal/dashboard"
	"moana/internal/store"
	"moana/internal/tz"
)

// Dashboard shows portfolio-style overview: balances, trends, charts, recent activity.
func (a *App) Dashboard(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	loc := tz.DisplayLocation(r)
	now := time.Now().UTC()

	data, err := dashboard.BuildPageData(ctx, a.Store, u.ID, loc, now, r.URL.Query().Get("period"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.renderShell(w, "dashboard.html", data, layoutShellMain("Dashboard", "dashboard", "dashboard-page", u))
}
