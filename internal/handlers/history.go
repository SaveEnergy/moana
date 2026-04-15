package handlers

import (
	"net/http"

	"moana/internal/historyview"
	"moana/internal/httperr"
	"moana/internal/store"
	"moana/internal/tz"
)

// History lists transactions with filters, search, sort, and date grouping.
func (a *App) History(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	loc := tz.DisplayLocation(r)

	data, err := historyview.BuildPage(ctx, a.Store, u.HouseholdID, loc, r.URL, r.URL.RequestURI())
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	a.historyRender(w, u, data)
}

func (a *App) historyRender(w http.ResponseWriter, u *store.User, data historyview.PageData) {
	a.renderShell(w, "history.html", data, layoutShell("History", "history", u))
}
